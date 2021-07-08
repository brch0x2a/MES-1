import requests


SERVERIP = "http://192.168.1.97:81"

YEAR = 2021
WEEK = 10

class Mecanic:
    def __init__(self, name):
        self.name = name
        self.planned = []
        self.actual = []

    def plannedToSaturation(self):
        saturation = [ t / 8 * 100 for t in self.planned[:-1]]
        saturation.append(self.planned[-1] / 48 * 100)
        return saturation

    def actualToSaturation(self):
        saturation = [ t / 8 * 100 for t in self.actual[:-1]]
        saturation.append(self.actual[-1] / 48 * 100)
        return saturation

class PM:
    def __init__(self, year, week):
        self.year  = year
        self.week = week
        self.Mecanics = []
        self.Frame  = []

    def process(self):
        params = dict(year=self.year, week=self.week)

        # --- Planeed ---
        route =  "/getPlannedWorkSaturarion"

        planned = getJsonRAW(route, params)

        for d in planned:
            name = d['Mecanic']['Fname'] #+ " " + d['Mecanic']['Lname']
            times  = d['Times']
            mecanic = Mecanic(name)
            mecanic.planned = times

            self.Mecanics.append(mecanic)

        # --- Actual ---
        route =  "/getActualWorkSaturarion"

        actual = getJsonRAW(route, params)

        i = 0
        for m in self.Mecanics:
            m.actual = actual[i]['Times']



    def getDataFrame(self):

        for i in range(len(self.Mecanics)):
            row = []
            row.append(self.Mecanics[i].name)
            row.append(self.Mecanics[i].planned[-1] / 48 * 100)
            row.append(self.Mecanics[i].actual[-1] / 48 * 100)

            # print("%s\t%f\t%f"%(row[0], row[1], row[2]))

            self.Frame.append(row)

        return self.Frame

    def plannedData(self):
        return [m.plannedToSaturation() for m in self.Mecanics]


class Ajob:
    """docstring for Autonomous."""

    def __init__(self, line, total, open, progress, closed):
        self.line = line
        self.total = total
        self.open = open
        self.progress = progress
        self.closed = closed

class AM:

    def __init__(self, init, end):
        self.init = init
        self.end = end
        self.Data = []
        self.Frame = []

    def process(self):
        params = dict(init=self.init, end=self.end)

        # --- Data ---
        route =  "/getAM_Stats"

        raw = getJsonRAW(route, params)

        for d in raw:
            line = d['Line']
            total = d['TotalJobs']
            open = d['OpenJobs']
            progress = d['InProgressJobs']
            closed = d['ClosedJobs']

            self.Data.append(Ajob(line, total, open, progress, closed))

    def getDataFrame(self):
        for i in range(len(self.Data)):
            row = []
            row.append(self.Data[i].line)
            # row.append(self.Data[i].total)
            row.append(self.Data[i].open)
            row.append(self.Data[i].progress)
            row.append(self.Data[i].closed)

            self.Frame.append(row)

        return self.Frame

class Allocation:

    def __init__(self, line, diff, allocated):
        self.line = line
        self.diff = diff
        self.allocated = allocated
        self.totalTime = abs(self.diff) + self.allocated

        if self.totalTime == 0:
            self.utilization = 0
        else:
            self.utilization = self.allocated / self.totalTime * 100

class OEELoss:
    def __init__(self, init, end):
        self.init = init
        self.end = end
        self.Data = []
        self.Frame = []

    def process(self):
        params = dict(start=self.init, end=self.end)

        # --- Data ---
        route =  "/getOEEMetaProjectionbyRange"

        raw = getJsonRAW(route, params)

        for d in raw:
            line = d['Line']
            diff = d['Diff']
            allocated = d['Allocated']

            self.Data.append(Allocation(line, diff, allocated))

    def getDataFrame(self):
        for i in range(len(self.Data)):
            row = []
            row.append(self.Data[i].line)
            row.append(self.Data[i].utilization)

            self.Frame.append(row)

        return self.Frame


def getJsonRAW(route, params):

    url =  SERVERIP + route

    res = requests.get(url, params=params)

    data = res.json()

    return data

#
# if __name__ == '__main__':

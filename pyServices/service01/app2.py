import pandas as pd
import numpy as np
import jinja2
from mail_controller import *
from engine import *


from flask import Flask, flash, request, redirect, url_for, render_template, send_from_directory
from flask_weasyprint import HTML, render_pdf

app = Flask(__name__)


WEEKDAYS = ["L", "K", "M", "J", "V", "S", "D", "Total"]

def processPMData(year, week):
    pm = PM(year, week)

    pm.process()


    pmFrameWeek = pd.DataFrame(pm.getDataFrame(), columns=["Mecanico", "Planeado", "Actual"], index=[m.name for m in pm.Mecanics])

    pmFramePlanned = pd.DataFrame(pm.plannedData(), columns=WEEKDAYS, index=[m.name for m in pm.Mecanics])

    # Plot
    ax = pmFrameWeek.plot.line(xlabel="Mecanicos", ylabel="%Utilizacion")
    fig = ax.get_figure()
    fig.savefig('static/images/plot.png')

    return pm


def processAMData(init, end):
    am = AM(init, end)

    am.process()

    am_frame = pd.DataFrame(am.getDataFrame(), columns=["Linea", "Abiertos", "Progresso", "Cerrado"], index=[m.line for m in am.Data])

    # Plot
    ax = am_frame.plot.line(xlabel="Linea", ylabel="Cantidad de Trabajos", color=["#ffd633", "#ff4dd2", "#4d4d4d"])
    fig = ax.get_figure()
    fig.savefig('static/images/am_plot.png')

    return am

def processOEEData(init, end):
    oee = OEELoss(init, end)

    oee.process()

    oee_frame = pd.DataFrame(oee.getDataFrame(), columns=["Linea", "Utilizacion"], index=[m.line for m in oee.Data])

    # Plot
    ax = oee_frame.plot.line(xlabel="Linea", ylabel="%Utilizacion", color=["#0000ff"])
    fig = ax.get_figure()
    fig.savefig('static/images/oee_plot.png')

    return oee


def startEngine():

    week = 12
    year = 2021
    init = "2021-03-15"
    end = "2021-03-22"

    pm = processPMData(year, week)
    am = processAMData(init, end)
    oee = processOEEData(init, end)

    # Template handling
    env = jinja2.Environment(loader=jinja2.FileSystemLoader(searchpath=''))
    template = env.get_template('template.html')
    html = template.render(pm=pm, am=am, oee=oee, title= str(year) + " WK " + str(week-1))


    # Write the HTML file
    with open('report.html', 'w') as f:
        f.write(html)

    mailEnngine =  MailEngine("Reporte ultilizacion", html)
    mailEnngine.start()

@app.route('/usageReport', methods=['GET', 'POST'])
def usageReport():

    week = 12
    year = 2021
    init = "2021-03-15"
    end = "2021-03-22"

    pm = processPMData(year, week)
    am = processAMData(init, end)
    oee = processOEEData(init, end)

    return render_template("template.html", pm=pm, am=am, oee=oee, title= str(year) + " WK " + str(week-1))




if __name__ == '__main__':
    # simple_table()
    #engineReport()
    # startEngine()
    app.run(debug=True, host="0.0.0.0", port="5001")

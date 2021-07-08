# simple_table.py
import pymysql
import os
import datetime

from fpdf import FPDF
from PIL import Image


from flask import Flask, flash, request, redirect, url_for, render_template, send_from_directory
from flask_weasyprint import HTML, render_pdf


db = pymysql.connect("192.168.1.97", "brch", "tk2718", "accesos")

dpi = 120 # note for calcs below that "pt" units are 1/72th of an inch

app = Flask(__name__)


class Access:

    def __init__(self, date_reg, in_charge, acces_type, person, temperature_epp, area, temperature, alert):
        self.date_reg = date_reg
        self.in_charge = in_charge
        self.acces_type = acces_type
        self.person = person
        self.temperature_epp = temperature_epp
        self.area = area
        self.temperature = temperature
        self.alert = alert


    def print(self):
        print("date:%s\t| in_charge:%s\t\t\t| acces_type:%s\t| person:%s \t| temperature:%s\t| alert:%s"%(
        self.date_reg, self.in_charge, self.acces_type, self.person, self.temperature, self.alert))


def get_reportRaw_by_date(init, end):

    cursor = db.cursor()

    sql = '''
    SELECT
        registro_fecha,
        Pr.persona_nombre,
        GETELEMENTO(registro_tipo_acceso_eid) AS registro_tipo_acceso,
        P.persona_nombre,
        CONCAT('images/registros/P10/',
                registro_id,
                '-E.jpg') AS foto_epp,
        CONCAT('images/registros/P10/',
                registro_id,
                '-A1.jpg') AS foto_area,
        R.registro_temperatura,
        R.registro_alerta
    FROM
        registros R
            INNER JOIN
        personas Pr ON R.registro_responsable_id = Pr.persona_id
            INNER JOIN
        personas P ON R.registro_persona_id = P.persona_id
    WHERE
        DATE(R.registro_fecha) BETWEEN %s AND %s
    '''

    cursor.execute(sql, (init, end))

    cursor.close()

    lines = cursor.fetchall()

    return lines


def simple_table(spacing=10):
    data = [['Fecha', 'Responsable', 'Acceso', 'Persona','Photo' 'Temperatura', 'Alerta'],
            ['2020-04-20 00:06:41', 'Driscoll', 'Entrada Principal 1', 'Kikolan', '/path', '37.60', '---'],

            ]

    pdf = FPDF()

    pdf.add_page()

    pdf.image('unilever.jpg', 10, 8, 33)
    pdf.set_font('Arial', 'B', 15)
    # Move to the right
    pdf.cell(80)
    # Title
    pdf.cell(30, 10, 'Reporte acceso COVID', 0, 0, 'C')
        # Line break


    pdf.set_font("Arial", size=8)

    col_width = pdf.w / 6
    row_height = pdf.font_size
    factor = 43
    relationImg = 32

    pdf.ln(row_height*spacing)


    for row in data:
        last_row = 0
        for item in row:
            # if last_row == 7:
            #     image = "house.jpg"
            #
            #     f = open(image, "rb")
            #
            #     im = Image.open(f)
            #
            #     pdf.cell(col_width, row_height*spacing, border=1)
            #
            #     pdf.image(image, col_width + 129, row_height+factor, relationImg, relationImg)
            #     print (image)
            #     f.close()
            #     factor += relationImg + 3.5
            #     pass

            pdf.cell(col_width, row_height*spacing,
                     txt=item, border=1)

            last_row += 1
            print(last_row+1)

        #
        pdf.ln(row_height*spacing)

    pdf.output('simple_table.pdf')


def buildTable(data):

    spacing = 10

    pdf = FPDF()

    pdf.add_page()

    pdf.image('unilever.jpg', 10, 8, 33)
    pdf.set_font('Arial', 'B', 15)
    # Move to the right
    pdf.cell(80)
    # Title
    pdf.cell(30, 10, 'Reporte acceso COVID', 0, 0, 'C')
        # Line break


    pdf.set_font("Arial", size=8)

    col_width = pdf.w / 6
    row_height = pdf.font_size
    factor = 43
    relationImg = 32

    pdf.ln(row_height*spacing)


    for row in data:
        last_row = 0
        for item in row:
            if last_row == 7:
                image = "house.jpg"

                f = open(image, "rb")

                im = Image.open(f)

                pdf.cell(col_width, row_height*spacing, border=1)

                pdf.image(image, col_width + 129, row_height+factor, relationImg, relationImg)
                print (image)
                f.close()
                factor += relationImg + 3.5
                pass

            pdf.cell(col_width, row_height*spacing,
                     txt=str(item), border=1)

            last_row += 1

            print(last_row+1)

        #
        pdf.ln(row_height*spacing)

    pdf.output('simple_table.pdf')



def engineReport():

    Report = []

    reportRaw = get_reportRaw_by_date("20200420", "20200420")

    for e in reportRaw:
        print("lenE: ", len(e))

        Report.append(Access(e[0], e[1], e[2], e[3], e[4], e[5]))


    for a in Report:
        a.print()

    buildTable(reportRaw)


def getReportBy(init, end):
    Report = []

    reportRaw = get_reportRaw_by_date(init, end)

    for e in reportRaw:
        # print("lenE: ", len(e))
        Report.append(Access(e[0], e[1], e[2], e[3], e[4], e[5], e[6], e[7]))


    return Report


@app.route('/report', methods=['GET', 'POST'])
def process_events_template():
    if request.method == 'POST':

        init = request.form['init']
        end = request.form['end']

        Report = getReportBy(init, end)

        # return render_template("reportDone.html", result = Report)
        html = render_template('reportDone.html', result = Report)
        return render_pdf(HTML(string=html))

    return render_template("reportFilter.html")




if __name__ == '__main__':
    # simple_table()
    #engineReport()
    app.run(debug=True, host="0.0.0.0")

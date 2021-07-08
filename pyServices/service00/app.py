import os
from flask import Flask, flash, request, redirect, url_for, render_template, send_from_directory
from werkzeug.utils import secure_filename
from enginePlanning import *
from engineEventMapping import *


UPLOAD_FOLDER = 'files/'
ALLOWED_EXTENSIONS = {"xlsx", "xlx"}

app = Flask(__name__)
app.config['UPLOAD_FOLDER'] = UPLOAD_FOLDER


def allowed_file(filename):
    return '.' in filename and \
           filename.rsplit('.', 1)[1].lower() in ALLOWED_EXTENSIONS



@app.route('/', methods=['GET', 'POST'])
def upload_file():
    if request.method == 'POST':
        # check if the post request has the file part
        if 'file' not in request.files:
            flash('No file part')
            print('No file part')
            return redirect(request.url)
        file = request.files['file']
        # if user does not select file, browser also
        # submit an empty part without filename
        if file.filename == '':
            flash('No selected file')
            print('No selected file')
            return redirect(request.url)
        if file and allowed_file(file.filename):
            filename = secure_filename(file.filename)
            file.save(os.path.join(app.config['UPLOAD_FOLDER'], filename))

            week = int(request.form['week'])
            year = int(request.form['year'])

            begin = beginDate(year, week)

            fileName = UPLOAD_FOLDER + filename


            notInDB = initPlanningEngine(fileName, begin)

            for e in notInDB:
                print("toRender: %s"%(e))

            print("file.save")
        return render_template("planDone.html", result = notInDB)

    return render_template('uploadTemplate.html')




@app.route('/get_planning_template', methods=['GET'])
def get_planning_template():

    return render_template('downloadTemplate.html')



@app.route('/get_events_template', methods=['GET', 'POST'])
def get_events_template():
    if request.method == 'POST':

        week = int(request.form['week'])
        year = int(request.form['year'])

        print("\nweek:%d\tyear:%d\n"%(week, year))

        start = time.time()

        begin = beginDate(year, week)

        endCalcDate = endDate(begin)

        print("\nweek:%d\tyear:%d\t%s --> %s\n"%(week, year, begin, endCalcDate))

        TEMP_FOLDER = "files/"

        file_name = "validacionOEE.xlsx"

        initEngineOEE(TEMP_FOLDER + file_name, begin)

        end = time.time()

        elapsed = (end - start) / 60
        print("Time elapsed: %f" % (elapsed))

        return send_from_directory(TEMP_FOLDER, file_name, as_attachment=True)

    return render_template("downloadValidationOEE.html")



@app.route('/process_events_template', methods=['GET', 'POST'])
def process_events_template():
    if request.method == 'POST':
        # check if the post request has the file part
        if 'file' not in request.files:
            flash('No file part')
            print('No file part')
            return redirect(request.url)
        file = request.files['file']
        # if user does not select file, browser also
        # submit an empty part without filename
        if file.filename == '':
            flash('No selected file')
            print('No selected file')
            return redirect(request.url)
        if file and allowed_file(file.filename):
            filename = secure_filename(file.filename)
            file.save(os.path.join(app.config['UPLOAD_FOLDER'], filename))

            week = int(request.form['week'])
            year = int(request.form['year'])

            begin = beginDate(year, week)

            file_name = UPLOAD_FOLDER + filename


            matchTest(file_name, begin)

            return render_template("processDone.html")


    return render_template('uploadEventTemplate.html')




if __name__ == '__main__':
    # app.secret_key = 'super secret key'
    app.run(debug=True, host="0.0.0.0")

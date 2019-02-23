import os
from flask import Flask, request, redirect, url_for
from flask import send_from_directory
from werkzeug.utils import secure_filename
import subprocess

MODELS = set(['la_muse.ckpt', 'rain_princess.ckpt', 'udnie.ckpt'])
UPLOAD_FOLDER = './upload'
RESULT_FOLDER = './result'
ALLOWED_EXTENSIONS = set(['png', 'jpg', 'jpeg'])
EVALUATE_COMMAND =   'python evaluate.py --allow-different-dimensions  --checkpoint ./models/udnie.ckpt --in-path ./upload/ --out-path ./result'

app = Flask(__name__)
app.config['MODELS'] = MODELS
app.config['RESULT_FOLDER'] = RESULT_FOLDER
app.config['UPLOAD_FOLDER'] = UPLOAD_FOLDER
app.config['EVALUATE_COMMAND'] = EVALUATE_COMMAND
def allowed_file(filename):
    return '.' in filename and \
           filename.rsplit('.', 1)[1].lower() in ALLOWED_EXTENSIONS

@app.route('/',methods=['GET'])


@app.route('/version',methods=['GET'])
def version():
	return "neural-art:v1.0.0"

@app.route('/', methods=['POST'])
def upload_file():
    if request.method == 'POST':
        # check if the post request has the file part
        if 'image' not in request.files:
            return 'No file'
        file = request.files['image']
        # if user does not select file, browser also
        # submit a empty part without filename
        if file.filename == '':
            return 'No selected file'
        if file and allowed_file(file.filename):
            filename = secure_filename(file.filename)
            file.save(os.path.join(app.config['UPLOAD_FOLDER'], filename))
        try:
        	result_success = subprocess.check_output(app.config['EVALUATE_COMMAND'], shell=True)
        except subprocess.CalledProcessError as e:
            return "An error occurred while evaluate"

        return send_from_directory(app.config['RESULT_FOLDER'],
                               filename)
    return '''Upload new File</title>'''

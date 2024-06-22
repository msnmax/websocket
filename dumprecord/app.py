from flask import Flask, render_template
import requests

app = Flask(__name__)


@app.route('/')
def home():
 

    return render_template('index.html')
 
@app.route('/dump')
def index():
 
        response = requests.get('http://localhost:8080/dump')
        if response.status_code == 200:
            data = response.text
            return data
        else:
            return f'Error fetching data: {response.status_code}'
   

if __name__ == '__main__':
    app.run(debug=True)
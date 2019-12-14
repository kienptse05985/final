from PIL import Image
from flask import Flask, request, jsonify
import numpy
from selenium import webdriver
from selenium.common.exceptions import NoAlertPresentException
from selenium.common.exceptions import TimeoutException
from selenium.common.exceptions import WebDriverException
from selenium.webdriver.common.alert import Alert
from tensorflow import keras
import time
import base64,os
from io import BytesIO

options = webdriver.ChromeOptions()
options.add_argument("--headless")
options.add_argument("--window-size=900,900")
options.add_argument("--no-sandbox")
options.add_argument("disable-infobars")
options.add_argument("--disable-extensions")
app = Flask(__name__)

model = keras.models.load_model('defactor_model.h5')


def send_response(response_obj):
    response = jsonify(response_obj)
    return response


# Resize 100x100
def resize(url):
    img = Image.open(url)
    image = img.convert("RGB")
    out_file = BytesIO()
    image.save(out_file, format='JPEG')
    im_data = out_file.getvalue()
    image_data = base64.b64encode(im_data)
    if not isinstance(image_data, str):
        image_data = image_data.decode()
    size = (img.size[0] / 9, img.size[1] / 9)
    img.thumbnail(size)
    # img.save(url)
    os.remove(url)
    # base64
    image = numpy.array(img.convert("RGB"))
    return image, image_data


def predict(image):
    image = image / 255
    image = image.reshape([-1, 100, 100, 3])
    prediction = model.predict_proba(image)
    return prediction[0][0] / (prediction[0][0] + prediction[0][1]), numpy.argmax(prediction[0])


# Route for handling the login page logic
@app.route('/scan', methods=['POST'])
def scan():
    result = {}
    if request.method == 'POST':
        body = request.json
        if body is None or 'url' not in body not in body:
            return send_response({
                'message': 'invalid request'
            }), 422
        if body['url'] == '':
            return send_response({
                'message': 'scan url can not be empty'
            }), 422
        scan_url = body['url']
        s_id = body['id']
        s_url = './image/' + s_id + '.png'

        browser = webdriver.Chrome(chrome_options=options, executable_path='./chromedriver')
        try:  # Alert
            browser.get(scan_url)
            #######
            Alert(browser).accept()  # accept the alert
            #######
            time.sleep(3)  # sleep to load the page for screenshot
            #######
            browser.save_screenshot(s_url)
            small, origin = resize(s_url)
            percent, pred = predict(small)
            result = {'message': 'ok', 'screen_shot': origin, 'code': 0, 'percentage': str(percent),
                      'prediction': str(pred)}
        except NoAlertPresentException:
            time.sleep(3)
            browser.save_screenshot(s_url)
            small, origin = resize(s_url)
            percent, pred = predict(small)
            result = {'message': 'ok', 'screen_shot': origin, 'code': 0, 'percentage': str(percent),
                      'prediction': str(pred)}
        except TimeoutException:  # timeout
            result = {'message': 'Timeout Error', 'code': 1}
            pass

        except WebDriverException:  # error page
            result = {'message': 'WebDriver Error', 'code': 1}
            pass
        browser.close()

    return send_response(result)


if __name__ == '__main__':
    app.run(host='localhost', port=5000, debug=True)

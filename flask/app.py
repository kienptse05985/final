from flask import Flask, render_template, redirect, url_for, request
from selenium import webdriver
from PIL import Image
import base64
from io import BytesIO  
from selenium.webdriver.common.alert import Alert
from selenium.common.exceptions import TimeoutException
from selenium.common.exceptions import NoAlertPresentException
from selenium.common.exceptions import WebDriverException
options = webdriver.ChromeOptions()
options.add_argument("--headless")
options.binary_location = "C:\Program Files (x86)\Viettel\SFive\Application\sfive.exe"
options.add_argument("--window-size=900,900")


app = Flask(__name__)
 
 
@app.route('/')
def welcome():
    return redirect('/capture')
url = './image/screen1.png'
#Resize 100x100
def resize(url):
    img = Image.open(url)
    size = (img.size[0]/9,  img.size[1]/9)
    img.thumbnail(size)
    img.save(url)
    #base64
    image = img.convert("RGB")
    out_file = BytesIO()
    image.save(out_file, format='JPEG')
    im_data = out_file.getvalue()
    image_data = base64.b64encode(im_data)
    if not isinstance(image_data, str):
        image_data = image_data.decode()
    return image_data

# Route for handling the login page logic
@app.route('/capture', methods=['GET', 'POST'])
def login():
    content = None
    if request.method == 'POST':
        loginurl = request.form['url']
        browser = webdriver.Chrome(chrome_options=options,executable_path=r"C:\Users\thanhnx23\Downloads\chromedriver_win32 (1)\chromedriver.exe")
        try:    # Alert 
            
                browser.get(loginurl)
                print loginurl
                #######
                Alert(browser).accept()#accept the alert
                print loginurl + '--Alert exists.... '
                #######
                time.sleep(5)#sleep to load the page for screenshot
                #######
                browser.save_screenshot(url)
                content = resize(url)
                print loginurl + '--Capture Alert'
                
        except NoAlertPresentException:             
                browser.save_screenshot(url)
                content = resize(url)
                print loginurl + '--Capture Without Alert'
                
        except TimeoutException: #timeout
                print loginurl + '--Error'
                content = 'Timeout Error'
                pass
            
        except WebDriverException: #error page
                print loginurl + '--WebDriverError'
                content = 'Web Error'
                pass
        browser.close()
        
    return render_template('index.html', output = content)


 
if __name__ == '__main__':
    app.run(host='localhost', port=5000, debug=True)

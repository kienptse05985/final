import csv
from PIL import Image
import time
from collections import defaultdict
from selenium import webdriver
from selenium.webdriver.common.alert import Alert
from selenium.common.exceptions import TimeoutException
from selenium.common.exceptions import NoAlertPresentException
from selenium.common.exceptions import WebDriverException
from selenium.webdriver.chrome.options import Options as ChromeOptions
options = ChromeOptions()
options.add_argument("--headless")
options.add_argument('window-size=900x900');
#host = "http://zonehmirrors.org" + raw_input("URL: ") 
#print host


urls = ['https://www.google.com']

#Resize 100x100
def resize(url):
    img = Image.open(url)
    size = (img.size[0]/9,  img.size[1]/9)
    img.thumbnail(size)
    img.save(url)
    print 'done'
    return url


#Screenshot
with webdriver.Chrome(executable_path=r'C:\Users\Acer\Downloads\chromedriver.exe',options=options) as driver:
    for index, url in enumerate(urls):
        output = {'screenshot': 'G:\\New folder\\image\\'+ str(index) +'_Deface7.png'}
        driver.set_page_load_timeout(10)
        try:    # Alert 
            
                driver.get(url)
                print url
                #######
                Alert(driver).accept()#accept the alert
                print str(index) + '--Alert exists.... '
                #######
                time.sleep(5)#sleep to load the page for screenshot
                #######
                driver.save_screenshot(output['screenshot'])
                resize(output['screenshot'])
                print str(index) + '--Capture Alert'
                
        except NoAlertPresentException:             
                driver.save_screenshot(output['screenshot'])
                resize(output['screenshot'])
                print str(index) + '--Capture Without Alert'
                
        except TimeoutException: #timeout
                print str(index) + '--Error'
                pass
            
        except WebDriverException: #error page
                print str(index) + '--WebDriverError'
                pass
        
    driver.close()
    
        


        
        
        


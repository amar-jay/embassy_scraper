from selenium import webdriver
# from selenium.webdriver.common.keys import Keys
from selenium.webdriver.common.by import By


URL = "https://www.google.com"
driver = webdriver.Chrome()

def main():
	driver.get(URL)
	driver.find_element(by=By.Attr, value="my-text")

	print("Hello, World!")

if __name__ == "__main__":
	main()
	driver.quit()
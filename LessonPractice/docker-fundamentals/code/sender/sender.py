import requests
import time
import logging
import sys

logging.basicConfig(level=logging.INFO, format='%(asctime)s - %(message)s',
                    datefmt='%Y-%m-%d %H:%M:%S', stream=sys.stdout)

while True:
    try:
        response = requests.get(f'http://receiver:8081')
        logging.info(
            f'Received response: {response.text}, status code: {response.status_code}')
        time.sleep(1)
    except requests.RequestException as e:
        logging.error("Error sending request to port 8081")
    finally:
        pass

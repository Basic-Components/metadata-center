FROM python:3.7
ADD pip.conf /etc/pip.conf
ADD ./requirements.txt /code/requirements.txt
WORKDIR /code
RUN pip install --upgrade pip
RUN pip install -r requirements.txt
ADD ./metadata_center /code/metadata_center
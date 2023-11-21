FROM python:3.9-bookworm

# Installing NodeJs to use nodemon
RUN apt-get update && apt-get install -y ca-certificates curl gnupg 
RUN mkdir -p /etc/apt/keyrings
RUN curl -fsSL https://deb.nodesource.com/gpgkey/nodesource-repo.gpg.key | gpg --dearmor -o /etc/apt/keyrings/nodesource.gpg
RUN echo "deb [signed-by=/etc/apt/keyrings/nodesource.gpg] https://deb.nodesource.com/node_20.x nodistro main" | tee /etc/apt/sources.list.d/nodesource.list
RUN apt-get update && apt-get install nodejs -y

RUN npm i -g nodemon
 
COPY requirements.txt /requirements.txt
RUN pip3 install -r /requirements.txt

WORKDIR /src 

CMD ["nodemon"] 
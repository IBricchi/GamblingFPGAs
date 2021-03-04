all: bin/meet_req run

bin/meet_req: requirements.txt
	pip3 install -r $<
	
	mkdir -p bin
	mkdir -p __pycache__
	touch bin/meet_req
	echo "*" > bin/.gitignore
	echo "*" > __pycache__/.gitignore

run:
	python3 app.py

clean:
	rm -rf bin
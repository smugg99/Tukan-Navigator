### This program is being hosted at https://navigator.smuggr.xyz

## Installation

Tested on Debian 12, you need to have GO (version 1.22.3) with all the dependencies from go.mod file installed and ready to use.

### To install some of the dependencies:

1. Install required packages:
	```
	sudo apt install -y wget tar nodejs npm
	```

2. Download Go 1.22.3 tarball:
	```
	wget https://go.dev/dl/go1.22.3.linux-amd64.tar.gz
	```

3. Extract the tarball to /usr/local:
	```
	sudo tar -C /usr/local -xzf go1.22.3.linux-amd64.tar.gz
	```

4. Set up the Go environment variables (Bash):
	Add this to ~/.bashrc
	```
	export PATH=$PATH:/usr/local/go/bin
	```
	and then
	```
	source ~/.bashrc
	```

### To setup Postgres (assuming you didn't change anything in the .env file):

1. Install postgresql:
	```
	sudo apt install -y postgresql postgresql-contrib
	```

2. Switch to postgres user and create the database structure:
	```
	sudo -i -u postgres
	```
	```
	psql
	```
	```
	CREATE USER smuggr WITH PASSWORD '1234567890';
	```
	```
	CREATE DATABASE net_work;
	```
	```
	GRANT ALL PRIVILEGES ON DATABASE net_work TO smuggr;
	```
	```
	psql -d net_work
	```
	```
	GRANT ALL PRIVILEGES ON SCHEMA public TO smuggr;
	GRANT ALL PRIVILEGES ON ALL TABLES IN SCHEMA public TO smuggr;
	GRANT ALL PRIVILEGES ON ALL SEQUENCES IN SCHEMA public TO smuggr;
	GRANT ALL PRIVILEGES ON ALL FUNCTIONS IN SCHEMA public TO smuggr;
	```
	```
	\q
	exit
	sudo systemctl restart postgresql
	```

### To compile/run the project:

1. Clone the repository to your local machine:
	```
	git clone https://github.com/smugg99/Tukan-Navigator.git
	```

2. Navigate to the project directory:
	```
	cd Tukan-Navigator
	```

3. Compile the whole project:
	```
	make
	```

4. Run the backend and frontend:
	```
	make run
	```

5. Open your web browser and visit `http://localhost:8000` (by default).

6. Check help of the makefile:
	```
	make help
	```

Toucan by Phil Laver from <a href="https://thenounproject.com/browse/icons/term/toucan/" target="_blank" title="Toucan Icons">Noun Project</a> (CC BY 3.0)

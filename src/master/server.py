#
# k8s-cms Master
# Flask Server
#

import settings
import contest

from flask import Flask
from healthcheck import HealthCheck
from sqlalchemy import create_engine

app = Flask(__name__)
# add api blueprints
app.register_blueprint(contest.api)
# setup healthcheck probe
health = HealthCheck(app, "/healthz")

# health check: check connection to db
@health.add_check
def check_db():
    engine = create_engine(settings.DB_CONNNECTION_STR)
    connection = engine.connect()
    results = connection.execute("SELECT 1")
    value = results[0][0]
    connection.close()

    if value == 1:
        return True, "db_ok"
    else:
        return False, "db_offline"

# root route
@app.route("/")
def route_root():
    return "k8s-cms master is up and running!"

if __name__ == "__main__":
    app.run(settings.SERVER_HOST, settings.LISTEN_PORT)

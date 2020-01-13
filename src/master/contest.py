#
# k8s-cms Master
# Contest API
#

from flask import request, Blueprint, jsonify

from settings import API_VERSION
from cms.db.util import SessionGen, Contest, get_contest_list

api = Blueprint("contest", __name__)
## routes
# api route lists contests on the cms system
@api.route(f"/api/v{API_VERSION}/contests", methods=["GET"])
def route_contests():
    with SessionGen() as session:
        contests = session.query(Contest).with_entities(Contest.id)
    contest_ids = [ contest.id for contest in contests ]
    return jsonify(contest_ids)


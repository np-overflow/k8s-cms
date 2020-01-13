#
# k8s-cms Master
# Contest API
#

from flask import abort, request, Blueprint, jsonify

from util import *
from settings import API_VERSION
from cms.db.util import SessionGen, Contest, get_contest_list

api = Blueprint("contest", __name__)

## utils
# get contest for the given id using the givne db session
def get_contest(session, contest_id):
    contest = session.query(Contest).get(contest_id)
    if contest is None: abort(404)

## mappings
# mapping from contest model to json representation
contest_mapping = {
    "name": "name",
    "description": "description",
    "allowed_localizations": "allowedLocalizations",
    "languages": "lanaguages",
    "submissions_download_allowed": "allowSubmissionsDownload",
    "allow_questions": "allowQuestions",
    "allow_password_authentication": "allowPasswordAuthentication",
    "ip_restriction": "enforceIPRestriction",
    "ip_autologin": "allowIPAutoLogin",
    "start": "startTime",
    "stop": "stopTime",
    "timezone": "timezone",
    "per_user_time": "maxContestPerUser",
    "max_submission_number": "maxSubmissionNum",
    "max_user_test_number": "maxTestsPerUser",
    "min_submission_interval": "minSubmissionInterval",
    "min_user_test_interval": "minUsertestInterval",
    "score_precision": "scorePrecision"
}

## routes
# api route lists contests on the cms systems by id with the given url param
# incl-names - whether to include names in json response
@api.route(f"/api/v{API_VERSION}/contests", methods=["GET"])
def route_contests():
    # get contests
    with SessionGen() as session:
        contests = get_contest_list()

    # parse filter params
    params = dict(request.args)
    has_include_names = parse_bool(params.get("incl-names", False))
    del params["incl-names"]
    if len(params) > 0: abort(400)

    if has_include_names:
        # return ids & names
        contests_meta = [ {"id": contest.id, "name": contest.name}
                         for contest in contests ]
        return jsonify(contests_meta)
    else: # return ids only
        contest_ids = [ contest.id for contest in contests ]
        return jsonify(contest_ids)

# api route to read, update contest metadata & delete contests
# GET - get the contest for the given contest id
# PATCH - update parameters on the contest for the given id 
# DELETE - delete the given contests for the given id
@api.route(f"/api/v{API_VERSION}/contest/<contest_id>", methods=["GET", "PATCH", "DELETE"])
def route_contest(contest_id):
    # check url params
    if contest_id is None: abort(400)

    with SessionGen() as session:
        if request.method == "GET":
            # get contest by id & map to json representation
            contest = get_contest(session, contest_id)
            contest_dict = map_dict(contest, contest_mapping)
            return jsonify(contest_dict)
        elif request.method == "PATCH":
            # update contest based on body params
            update_params = dict(request.json)
            contest = get_contest(session, contest_id)
            contest = map_obj(contest_id, update_params, reverse_mapping(contest_mapping))
            session.commit()
        elif request.method == "DELETE":
            # delete contest based oin body params
            contest = get_contest(session, contest_id)
            session.delete(contest)
        else:
            abort(400) # unimplemented http method

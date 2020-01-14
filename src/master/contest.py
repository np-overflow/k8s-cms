#
# k8s-cms Master
# Contest API
#

import os
import tarfile
import tempfile
import subprocess
from shutil import rmtree
from flask import abort, request, Blueprint, jsonify

from util import *
from auth import authenticate
from settings import API_VERSION
from cms.db.util import SessionGen, Contest, get_contest_list
from cmscontrib.ImportContest import ContestImporter
from cmscontrib.ImportUser import UserImporter, choose_loader
from cmscontrib.loaders import choose_loader, build_epilog

api = Blueprint("contest", __name__)

## contests utils
# unpack formdata contest data payload for request
# returns path of the unpacked data
def unpack_data(request):
    # check contest data payload present
    if not "file" in request.files: abort(400)
    file = request.files["file"]
    # write file data to temp file & extract contents to dir
    with tempfile.NamedTemporaryFile("r") as f:
        file.save(f.name)
        with tarfile.open(f.name,  "r:gz") as tar:
            data_dir = tempfile.mkdtemp()
            tar.extractall(data_dir)
    return data_dir

# import the given contest using data from data_dir
# is_update - whether contest data import is used to update an existing contest
# returns the id of the imported contest
def import_contest(data_dir, for_update=False):
    # attempt autodetect contest loader
    def error_callback(message):
        abort(400, message)
    loader_class = choose_loader(
        None, # autodetect loader
        data_dir,
        error_callback # called when unable to autodetect
    )

    # import users in contest
    get_loader = (lambda p: loader_class)
    UserImporter(data_dir,
                 contest_id=None,
                 loader_class=loader_class).do_import_all(data_dir, get_loader)

    # import contest
    contest_id =  ContestImporter(data_dir,
                                  yes=True, # dont prompt on delete
                                  zero_time=False,
                                  import_tasks=True,
                                  no_statements=False,
                                  update_contest=for_update,
                                  update_tasks=for_update,
                                  delete_stale_participations=True,
                                  loader_class=loader_class).do_import()

    # check if create import caused conflict in db as contest already exists
    if contest_id == False and not for_update: abort(409)

    return contest_id


## mappings
# mapping from contest model to json representation
contest_mapping = [
    ("name", "name"),
    ("description", "description"),
    ("allowed_localizations", "allowedLocalizations"),
    ("languages", "languages"),
    ("submissions_download_allowed", "allowSubmissionsDownload"),
    ("allow_questions", "allowQuestions"),
    ("allow_password_authentication", "allowPasswordAuthentication"),
    ("ip_restriction", "enforceIPRestriction"),
    ("ip_autologin", "allowIPAutoLogin"),
    ("start", "startTime"),
    ("stop", "stopTime"),
    ("timezone", "timezone"),
    ("per_user_time", "maxContestPerUser"),
    ("max_submission_number", "maxSubmissionNum"),
    ("max_user_test_number", "maxTestsPerUser"),
    ("min_submission_interval", "minSubmissionInterval"),
    ("min_user_test_interval", "minUsertestInterval"),
    ("score_precision", "scorePrecision"),
]

## routes
# api route lists contests on the cms systems by id with the given url param
# incl-names - whether to include names in json response
@api.route(f"/api/v{API_VERSION}/contests", methods=["GET"])
@authenticate(kind="access")
def route_contests():
    # parse filter params
    params = dict(request.args)
    # incl-names filter param - include names in response
    has_include_names = parse_bool(params.get("incl-names", False))
    if has_include_names: del params["incl-names"]
    # check for unknown filter params
    if len(params) > 0: abort(400)

    with SessionGen() as session:
        if has_include_names:
            # return ids & names
            contests = session.query(Contest).with_entities(Contest.id, Contest.name)
            contests_meta = [ {"id": contest.id, "name": contest.name}
                             for contest in contests ]
            return jsonify(contests_meta)
        else: # return ids only
            contests = session.query(Contest).with_entities(Contest.id)
            contest_ids = [ contest[0] for contest in contests ]
            return jsonify(contest_ids)

# api route to read, update contest metadata & delete contests
# GET - get the contest for the given contest id
# PATCH - update parameters on the contest for the given id 
# DELETE - delete the given contests for the given id
@api.route(f"/api/v{API_VERSION}/contest/<int:contest_id>", methods=["GET", "PATCH", "DELETE"])
@authenticate(kind="access")
def route_contest(contest_id):
    # check url params
    if contest_id is None: abort(400)

    with SessionGen() as session:
        if request.method == "GET":
            # get contest by id & map to json representation
            contest = get(Contest, session, contest_id)
            contest_dict = map_dict(contest, contest_mapping)
            return jsonify(contest_dict)
        elif request.method == "PATCH":
            # update contest based on body params
            update_params = dict(request.json)
            contest = get(Contest, session, contest_id)
            contest = map_obj(contest, update_params, reverse_mapping(contest_mapping))
            session.commit()
            return jsonify({"success": True})
        elif request.method == "DELETE":
            # delete contest based on body params
            # TODO: kill ProxyService for contest
            contest = get(Contest, session, contest_id)
            session.delete(contest)
            session.commit()
            return jsonify({"success": True})
        else:
            abort(404) # unimplemented http method

# api route to import, update contests based on contest data in the formats supported by CMS
# request should a formdata file: gzipped tar archive of contest data 
# POST - create a new contest for the given contest data
# PATCH - update the contest for the given contest_id with the given contest data
@api.route(f"/api/v{API_VERSION}/contest/import", methods=["POST"])
@api.route(f"/api/v{API_VERSION}/contest/import/<int:contest_id>", methods=["PATCH"])
@authenticate(kind="access")
def route_import_contest(contest_id=None):
    # unpack contest data 
    data_dir = unpack_data(request)

    if request.method == "POST":
        # TODO: spawn ProxyService for contest
        contest_id = import_contest(data_dir)
    elif request.method == "PATCH":
        # verify that a contest with id exists
        with SessionGen() as session:
            contest = session.query(Contest).get(contest_id)
            if contest is None: abort(404)

        contest_id = import_contest(data_dir, for_update=True)

    # cleanup
    rmtree(data_dir)

    return jsonify({"id": contest_id})

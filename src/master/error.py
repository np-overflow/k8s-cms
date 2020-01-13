#
# k8s-cms master
# Error Handlers
#

import settings
import contest

from flask import Blueprint, jsonify

handlers = Blueprint("error", __name__)

# error handler for bad requests
@handlers.app_errorhandler(400)
def handle_bad_request(error):
    return jsonify({"error": "Bad Request", "message": error.description}), 400

# error handler for not found
@handlers.app_errorhandler(404)
def handle_bad_request(error):
    return jsonify({"error": "Not Found", "message": error.description}), 404

# error handler for not found
@handlers.app_errorhandler(409)
def handle_bad_request(error):
    return jsonify({"error": "Conflict", "message": error.description}), 409

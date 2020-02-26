#
# k8s-cms Master
# Authentication API
#

import os
import jwt
from functools import wraps
from util import map_obj, map_dict, reverse_mapping
from datetime import datetime, timedelta
from base64 import b64encode, b64decode
from flask import abort, request, Blueprint, jsonify

from settings import API_VERSION, JWT_KEY
from cms.db import SessionGen, Admin
from cmscommon.crypto import validate_password


api = Blueprint("auth", __name__)
## models
# represents a token that is use to authenticate
class Token:
    JWT_ISSUER = "k8s-cms/master"
    JWT_AUDIENCE = "k8s-cms/master"

    # mapping between jwt claims and token field
    mapping = {
        ("subject", "sub"),
        ("issue_on", "iat"),
        ("expires", "exp"),
    }

    # create a token of the given kind
    # user_id - unique identifier of the user that created the token
    def __init__(self, kind=None, user_id=None):
        self.issue_on = datetime.utcnow()
        self.subject = f"k8s-cms/token/{kind}/{user_id}"

        if kind == "access":
            self.expires = self.issue_on + timedelta(minutes=5)
        elif kind == "refresh":
            self.expires = self.issue_on + timedelta(days=14)
        elif kind is None:
            self.expires = None
        else:
            raise ValueError(f"Unsupported token kind: {kind}")


    # convert this object as a JWT token
    # returns this object in JWT token representation
    def to_jwt(self):
        # build payload for jwt
        payload = map_dict(self, Token.mapping)
        jwt_bytes = jwt.encode(payload, JWT_KEY, algorithm="HS256")
        return b64encode(jwt_bytes).decode()

    # create token using the given JWT token
    @classmethod
    def from_jwt(cls, jwt_token):
        # extract jwt
        jwt_bytes = b64decode(jwt_token)
        payload = jwt.decode(jwt_bytes, JWT_KEY, algorithms='HS256')
        # map fields
        token = cls()
        map_obj(token, payload, reverse_mapping(Token.mapping))
        # convert data types
        token.issue_on = datetime.fromtimestamp(token.issue_on)
        token.expires = datetime.fromtimestamp(token.expires)

        return token

    ## properties
    @property
    def kind(self):
        _, _, kind, _, _ = self.subject.split("/")
        return kind

    @property
    def user_id(self):
        _, _, _, _, user_id = self.subject.split("/")
        return user_id

## utils

# perform login for the given credentials
# returns id of admin if login is success
def perform_login(username, password):
    with SessionGen() as session:
        # attempt to get admin for username & check if enabled
        admin = session.query(Admin).filter_by(username=username).one_or_none()
        if admin is None or not admin.enabled: abort(401)

        # attempt to validate password
        try:
            allowed = validate_password(admin.authentication, password)
        except ValueError:
            allowed = False
        if not allowed: abort(401)

        return admin.id

# kind - kind of token to validate "access, "refresh" or "any"
def verify_token(jwt_token, kind="access"):
    # attempt to parse JWT token
    try:
        token = Token.from_jwt(jwt_token)
    except Exception as e:
        return False

    # verify token fields
    if kind != "any" and token.kind != kind: return False
    if token.expires < datetime.utcnow(): return False
    if token.issue_on > datetime.utcnow(): return False
    # check if user exists
    with SessionGen() as session:
        admin = session.query(Admin).get(token.user_id)
        if admin is None: return False

    return True

# decorator to verify authentication before executing the given request fn
# kind - kind of token to validate "access, "refresh" or "any"
def authenticate(kind="access"):
    def auth_decorator(request_fn):
        @wraps(request_fn)
        def auth_fn(*args, **kwargs):
            # check auth header
            auth_header = request.headers.get("Authorization", "")
            if not "Bearer" in auth_header: abort(401)
            # extract & verify jwt token
            _, jwt_token = auth_header.split()
            if not verify_token(jwt_token, kind): abort(401)
            # run request fn
            return request_fn(*args, **kwargs)
        return auth_fn
    return auth_decorator

## routes
# api route to perform login with login credentials
# POST - perform login with login credentials
# responses with a refresh token 
@api.route(f"/api/v{API_VERSION}/auth/login", methods=["POST"])
def route_login():
    # check if can parse body as json
    credentials = request.json
    if credentials is None: abort(400)

    # perform login & generate refresh token
    uid = perform_login(credentials["username"], credentials["password"])
    admin_id = f"admin/{uid}"
    jwt_token = Token("refresh", admin_id).to_jwt()

    return jsonify({"refreshToken": jwt_token })

# api route to refresh access token with refresh token
# POST - refresh access token with refresh token
# responses with a access token
@api.route(f"/api/v{API_VERSION}/auth/refresh", methods=["GET"])
@authenticate(kind="refresh")
def route_refresh():
    # extract user id from refresh token
    auth_header = request.headers.get("Authorization", "")
    _, refresh_jwt_token = auth_header.split()
    refresh_token = Token.from_jwt(refresh_jwt_token)

    # generate access token
    admin_id = f"admin/{refresh_token.user_id}"
    access_jwt_token = Token("access", admin_id).to_jwt()

    return jsonify({"accessToken": access_jwt_token })

# api route to verify if token is valid
@api.route(f"/api/v{API_VERSION}/auth/check", methods=["GET"])
@authenticate(kind="any")
def route_check():
    return jsonify({"success": True})


#
# k8s-cms master
# settings
#

import os
import json
import string
import random

# cms config
CMS_CONFIG_PATH = os.environ.get("CMS_CONFIG", "/cms/config/cms.conf")
with open(CMS_CONFIG_PATH, "r") as f:
    cms_config = json.load(f)

# api settings
API_VERSION = 0

# server settings
SERVER_HOST = os.environ.get("KCMS_MASTER_HOST", cms_config["master_listen_address"])
LISTEN_PORT = os.environ.get("KCMS_MASTER_PORT", cms_config["master_listen_port"])

# db settings
DB_CONNNECTION_STR = cms_config["database"]

# jwt settings
keyspace = string.digits + string.ascii_letters
default_jwt_key = "".join(random.choices(keyspace, k=32))
JWT_KEY = os.environ.get("KCMS_MASTER_JWT_KEY", default_jwt_key)

#
# k8s-cms master
# settings
#

import os
import json

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
JWT_KEY = os.environ.get("KCMS_MASTER_JWT_KEY", cms_config["master_jwt_key"])

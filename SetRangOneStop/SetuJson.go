package SetRangeOneStop

const PermissionJson = `{
    "policyType": 0,
    "description": "STARROCKS.PARMS.DESCRIPTION",
    "isEnabled": true,
    "isNormal": false,
    "createdBy": "",
    "updatedBy": "",
    "service": "STARROCKS.SERVICE.NAME",
    "name": "STARROCKS.PARMS.POLICYNAME",
    "isAuditEnabled": true,
    "resources": {
        "database": {"values": STARROCKS.PARMS.DATABASE,"isExcludes": false,"isRecursive": false},
        "column": {"values": STARROCKS.PARMS.COLUMN,"isExcludes": false,"isRecursive": false},
        "catalog": {"values": STARROCKS.PARMS.CATALOG,"isExcludes": false,"isRecursive": false},
        "table": {"values": STARROCKS.PARMS.TABLE,"isExcludes": false,"isRecursive": false}
    },
    "policyItems": [
        {"users": STARROCKS.PARMS.USER,"accesses": STARROCKS.PARMS.PERMISSION,
            "groups": [""],
            "roles": null,
            "conditions": null,
            "delegateAdmin": false
        }
    ],
    "options": {"POLICY_VALIDITY_SCHEDULES": ""},
    "validitySchedules": null,
    "policyLabels": ["STARROCKS.PARMS.POLICYLABELS"],
    "zone_name": ""
}`

const PerJson = `{
    "policyType": 0,
    "description": "STARROCKS.PARMS.DESCRIPTION",
    "isEnabled": true,
    "isNormal": false,
    "createdBy": "",
    "updatedBy": "",
    "service": "STARROCKS.SERVICE.NAME",
    "name": "STARROCKS.PARMS.POLICYNAME",
    "isAuditEnabled": true,
    "resources": STARROCKS.PARMS.RESOURCES,
    "policyItems": [
        {"users": STARROCKS.PARMS.USER,"accesses": STARROCKS.PARMS.PERMISSION,
            "groups": [""],
            "roles": null,
            "conditions": null,
            "delegateAdmin": false
        }
    ],
    "options": {"POLICY_VALIDITY_SCHEDULES": ""},
    "validitySchedules": null,
    "policyLabels": ["STARROCKS.PARMS.POLICYLABELS"],
    "zone_name": ""
}`

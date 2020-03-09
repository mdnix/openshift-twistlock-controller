#!/usr/bin/env bash

oc create sa twistlock-cluster-reader -n rch-twistlock-sync-tst
oc adm policy add-cluster-role-to-user cluster-reader -z twistlock-cluster-reader -n rch-twistlock-sync-tst
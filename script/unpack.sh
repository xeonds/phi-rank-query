#!/bin/bash
APK=$1
python3 gameInfo.py $APK && \
python3 resource.py $APK && \
mv IllustrationLowRes ../web/public/assets/illustrations && \
mv *.csv ../web/public/
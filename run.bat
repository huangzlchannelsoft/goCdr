cdrAlarm.exe -as client -svrAddr 0.0.0.0:9080 -cdrPath /home/umg/ATS4/cdr
#cdrAlarm.exe -as server -svrAddr 0.0.0.0:9080 
              -alarmUri http://127.0.0.1:9070/sendAlarm 
              -pushGateWayUri http://10.130.41.226:9091
#cdrAlarm.exe -as singleton 
              -alarmUri http://127.0.0.1:9070/hooks/cdrAlarm
              -pushGateWayUri http://10.130.41.226:9091

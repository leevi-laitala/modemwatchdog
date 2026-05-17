from flask import Flask, jsonify, request

app = Flask(__name__)

powerOn = False
hasInternet = True

@app.route("/gen_204", methods = ["HEAD"])
def pulse():
    global powerOn
    global hasInternet

    if not hasInternet:
        return "", 400

    return "", 204

@app.route("/nointernet", methods = ["POST"])
def nointernet():
    global hasInternet

    hasInternet = False
    return "", 200

@app.route("/internet", methods = ["POST"])
def internet():
    global hasInternet

    hasInternet = True
    return "", 200

@app.route("/turnoff", methods = ["POST"])
def turnOff():
    global powerOn

    statusWas = powerOn
    powerOn = False
    return jsonify({"msg": "Turned off",
                    "currentStatus": "on" if powerOn else "off",
                    "statusWas": "on" if statusWas else "off"}), 200


@app.route("/turnon", methods = ["POST"])
def turnOn():
    global powerOn

    statusWas = powerOn
    powerOn = True
    return jsonify({"msg": "Turned on",
                    "currentStatus": "on" if powerOn else "off",
                    "statusWas": "on" if statusWas else "off"}), 200


@app.route("/status", methods = ["GET"])
def powerStatus():
    global powerOn

    return jsonify({"currentStatus": "on" if powerOn else "off"}), 200

app.run()

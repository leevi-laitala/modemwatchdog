from flask import Flask, jsonify, request

app = Flask(__name__)

powerOn = False

attempt = 0
maxAttempt = 3
def failSometimes() -> bool:
    global attempt
    global maxAttempt

    attempt += 1
    attempt = attempt % maxAttempt

    return attempt != 0

@app.route("/gen_204", methods = ["HEAD"])
def pulse():
    global powerOn

    if not failSometimes():
        return jsonify({"msg": "fail"}), 403

    return "", 204

@app.route("/api/turnoff", methods = ["POST"])
def turnOff():
    global powerOn

    if failSometimes():
        return jsonify({"msg": "fail"}), 403

    statusWas = powerOn
    powerOn = False
    return jsonify({"msg": "Turned off",
                    "currentStatus": "on" if powerOn else "off",
                    "statusWas": "on" if statusWas else "off"}), 200


@app.route("/api/turnon", methods = ["POST"])
def turnOn():
    global powerOn

    if failSometimes():
        return jsonify({"msg": "fail"}), 403

    statusWas = powerOn
    powerOn = True
    return jsonify({"msg": "Turned on",
                    "currentStatus": "on" if powerOn else "off",
                    "statusWas": "on" if statusWas else "off"}), 200


@app.route("/api/status", methods = ["GET"])
def powerStatus():
    global powerOn

    if failSometimes():
        return jsonify({"msg": "fail"}), 403

    return jsonify({"currentStatus": "on" if powerOn else "off"}), 200

app.run()

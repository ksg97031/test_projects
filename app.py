import sys
import hashlib
from flask import Flask, render_template, request, session, jsonify

app = Flask(__name__)
app.secret_key = "6fc40225787ef901de64c056a729103d1f0b1e5fdda9a1db" # random generate

@app.route('/')
def index():
    t = request.args.get('input')
    return t.format(test=123)


if __name__ == "__main__":
    port = 80
    if len(sys.argv) > 1:
        port = int(sys.argv[1])

    app.run(host='0.0.0.0', port=port)


import sys
import json

import requirements


def process(request):
    return {
        'files': [
            {
                'name': file_.get('name'),
                'dependencies': [
                    {
                        'name': req.name,
                        'specifiers': [
                            {
                                'operator': spec[0],
                                'version': spec[1],
                            }
                            for spec in req.specs
                        ],
                        'extras': {
                            key: getattr(req, key)
                            for key in (
                                'line',
                                'editable',
                                'vcs',
                                'revision',
                                'uri',
                                'path',
                                'extras'
                            )
                            if getattr(req, key)
                        },
                    }
                    for req in requirements.parse(file_.get('content', ''))
                ],
            }
            for file_ in request['files']
        ],
    }


if __name__ == '__main__':
    request = json.load(sys.stdin)
    response = process(request)
    json.dump(response, sys.stdout)

import os, sys, unittest

os.environ['EVIDENTLY_PROJECT_NAME'] = 'retaildemostore'

sys.argv += ['discover', os.path.dirname(sys.argv[0]), 'test_*.py']

unittest.main(module=None)
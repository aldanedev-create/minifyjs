#!/usr/bin/env python3
"""Install a wheel in a temporary venv and verify its shipped CLI and API."""
import argparse
import os
from pathlib import Path
import subprocess
import tempfile
import venv


def main():
    parser = argparse.ArgumentParser(description=__doc__)
    parser.add_argument('wheel', type=Path)
    args = parser.parse_args()
    wheel = args.wheel.resolve(strict=True)
    version = wheel.name.split('-')[1]
    smoke = Path(__file__).with_name('wheel_smoke.py').resolve()
    with tempfile.TemporaryDirectory(prefix='minifyjs-wheel-') as tmp:
        root = Path(tmp)
        env_dir = root / 'venv'
        venv.EnvBuilder(with_pip=True).create(env_dir)
        python = env_dir / ('Scripts/python.exe' if os.name == 'nt' else 'bin/python')
        env = os.environ.copy()
        for key in ('PYTHONPATH', 'MINIFYJS_BINARY'):
            env.pop(key, None)
        subprocess.run([str(python), '-m', 'pip', 'install', '--no-index', '--no-deps', str(wheel)], cwd=root, env=env, check=True)
        subprocess.run([str(python), '-I', str(smoke), version], cwd=root, env=env, check=True)


if __name__ == '__main__':
    main()

"""Exercise an installed wheel from outside the source checkout."""
import json
import subprocess
import sys
from pathlib import Path

import minifyjs

expected = sys.argv[1]
assert minifyjs.__version__ == expected, minifyjs.__version__
assert subprocess.check_output([sys.executable, '-m', 'minifyjs', '--version'], text=True).strip().startswith(f"minifyjs {expected} ")
cli = Path(sys.executable).parent / ('minifyjs.exe' if sys.platform == 'win32' else 'minifyjs')
assert cli.is_file(), 'console entry point missing'
assert subprocess.check_output([str(cli), '--version'], text=True).strip().startswith(f"minifyjs {expected} ")
source = 'export function calculate(longParameter) { return longParameter + 1; }'
result = minifyjs.minify(source, format='esm', mangle=True)
assert 'longParameter' not in result.code
assert len(result.code.encode()) < len(source.encode())
Path('input.js').write_text(source)
subprocess.run([str(cli), 'input.js', '-o', 'output.js', '--mangle'], check=True)
assert Path('output.js').is_file()
project = Path('project').resolve()
project.mkdir()
(project / 'main.js').write_text('import {used} from "./exports.js"; export const answer = used; export const load = () => import("./lazy.js"); if (DEBUG) console.log("DEBUG_SENTINEL"); debugger;')
(project / 'exports.js').write_text('export {used, unused} from "./values.js";')
(project / 'values.js').write_text('export const used = 42; export const unused = "UNUSED_SENTINEL";')
(project / 'lazy.js').write_text('export const value = "lazy-loaded";')
result = minifyjs.bundle(['main.js'], working_dir=str(project), outdir='dist',
    splitting=True, format='esm', target='es2020', entry_names='[name]-[hash]',
    chunk_names='chunks/[name]-[hash]', define={'DEBUG':'false'}, drop=['debugger'],
    sourcemap='external', metafile='meta.json')
paths = [Path(f['path']) for f in result.output_files]
assert all(p.is_file() for p in paths)
assert any(p.parent.name == 'chunks' and p.suffix == '.js' for p in paths)
assert any(p.suffix == '.map' for p in paths)
assert result.minified_bytes == sum(p.stat().st_size for p in paths)
assert result.metafile == json.loads((project / 'meta.json').read_text())
code = '\n'.join(p.read_text() for p in paths if p.suffix == '.js')
assert 'UNUSED_SENTINEL' not in code and 'DEBUG_SENTINEL' not in code and 'debugger' not in code
package = project / 'node_modules' / 'fixture'
package.mkdir(parents=True)
(package / 'package.json').write_text('{"name":"fixture","main":"index.js"}')
(package / 'index.js').write_text('exports.value = 123;')
(project / 'package.js').write_text('import {value} from "fixture"; console.log(value);')
bundled = minifyjs.bundle(['package.js'], working_dir=str(project), outfile='bundled.js')
assert bundled.code and bundled.minified_bytes == len(bundled.code.encode())
assert 'fixture' not in bundled.code
external = minifyjs.bundle(['package.js'], working_dir=str(project), outfile='external.js', packages='external')
assert 'fixture' in external.code
print('Installed wheel smoke checks passed for', expected)

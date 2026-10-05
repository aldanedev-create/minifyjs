"""Validate MinifyJS bundles against a pinned Teloce checkout in Chromium."""
import argparse
from functools import partial
from http.server import SimpleHTTPRequestHandler, ThreadingHTTPServer
from pathlib import Path
import subprocess
import tempfile
from threading import Thread

from minifyjs import bundle
from teloce.build.builder import Builder
from teloce.router.compiler import RouterCompiler
from teloce.router.generator import RouterGenerator


def prepare(root, teloce, esbuild):
    for shared in (False, True):
        project = root / str(shared).lower()
        project.mkdir()
        (project / 'App.vel').write_text((teloce / 'examples/flask/static/js/App.vel').read_text())
        generated = project / 'generated'
        result = Builder({'dev': True, 'minify': False, 'source_maps': False,
                          'shared_runtime': shared, 'spa': False}).build(project, generated)
        assert result['failed'] == 0, result
        (generated / 'Lazy.js').write_text('export const value = "lazy-loaded";')
        config = RouterCompiler().compile({'routes': [
            {'path': '/', 'component': 'Home'},
            {'path': '/users/:id', 'component': 'User'}]})
        (generated / 'Router.js').write_text('const Home = "Home", User = "User";\n' + RouterGenerator().generate(config))
        (generated / 'entry.js').write_text('import {mount} from "./App.js"; import router from "./Router.js"; mount("#app"); window.testRouter = router; window.loadChunk = () => import("./Lazy.js");')
        dist = project / 'dist'
        result = bundle(['entry.js'], working_dir=str(generated), outdir=str(dist),
                        splitting=True, format='esm', target='es2020', mangle=True,
                        chunk_names='chunks/[name]-[hash]')
        assert any('/chunks/' in f['path'].replace('\\', '/') for f in result.output_files)
        # Compare identical settings and input with the independent Node CLI.
        node_dist = project / 'node-dist'
        subprocess.run([str(esbuild),
            'entry.js', '--bundle', '--minify', '--format=esm', '--target=es2020',
            '--charset=utf8', '--splitting', '--chunk-names=chunks/[name]-[hash]',
            '--outdir=' + str(node_dist)], cwd=generated, check=True)
        min_files = {p.relative_to(dist).as_posix(): p.read_bytes() for p in dist.rglob('*.js')}
        node_files = {p.relative_to(node_dist).as_posix(): p.read_bytes() for p in node_dist.rglob('*.js')}
        assert min_files == node_files, 'Node esbuild output differs'
        (dist / 'index.html').write_text('<div id="app"></div><script type="module" src="entry.js"></script>')


def main():
    parser = argparse.ArgumentParser(description=__doc__)
    parser.add_argument('teloce', type=Path)
    parser.add_argument('--prepare-only', action='store_true')
    parser.add_argument('--esbuild', type=Path, required=True)
    args = parser.parse_args()
    with tempfile.TemporaryDirectory(prefix='minifyjs-teloce-') as tmp:
        root = Path(tmp)
        prepare(root, args.teloce.resolve(), args.esbuild.resolve())
        if args.prepare_only:
            print('Teloce fixtures and Node esbuild byte parity passed')
            return
        from playwright.sync_api import sync_playwright
        server = ThreadingHTTPServer(('127.0.0.1', 0), partial(SimpleHTTPRequestHandler, directory=str(root)))
        Thread(target=server.serve_forever, daemon=True).start()
        try:
            with sync_playwright() as pw:
                browser = pw.chromium.launch()
                try:
                    for shared in ('false', 'true'):
                        page = browser.new_page()
                        errors = []
                        page.on('pageerror', lambda error: errors.append(str(error)))
                        page.goto(f'http://127.0.0.1:{server.server_port}/{shared}/dist/index.html')
                        page.get_by_role('heading', name='Flask dashboard').wait_for()
                        page.get_by_role('button', name='Sign out', exact=True).click()
                        page.get_by_role('button', name='Sign in', exact=True).wait_for()
                        assert page.locator('.status').count() == 0
                        page.get_by_role('button', name='Sign in', exact=True).click()
                        page.locator('.status').wait_for()
                        state = page.evaluate('async () => { await window.testRouter.push("/users/42?tab=posts"); return {id: window.testRouter.params().id, tab: window.testRouter.query().tab, component: window.testRouter.currentRoute().component}; }')
                        assert state == {'id':'42', 'tab':'posts', 'component':'User'}, state
                        assert page.evaluate('async () => (await window.loadChunk()).value') == 'lazy-loaded'
                        assert not errors, errors
                        page.close()
                    print('Teloce rendering, reactivity, routing, and lazy chunks passed in both runtime modes')
                finally:
                    browser.close()
        finally:
            server.shutdown()
            server.server_close()


if __name__ == '__main__':
    main()

import * as esbuild from 'esbuild';

await esbuild.build({
  entryPoints: ['ts/main.page.ts'],
  bundle: true,
  minify: false,
  outfile: 'templates/js/main.page.js',
});

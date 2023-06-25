import * as esbuild from 'esbuild';

await esbuild.build({
  entryPoints: ['ts/signin.page.ts'],
  bundle: true,
  minify: false,
  outfile: 'templates/js/signin.page.js',
});

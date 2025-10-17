module.exports = {
  branches: ['main'],
  tagFormat: 'v${version}',
  plugins: [
    '@semantic-release/commit-analyzer',
    '@semantic-release/release-notes-generator',
    ['@semantic-release/exec', {
      prepareCmd: './build.sh ${nextRelease.version}',
    }],
    ['@semantic-release/github', {
      assets: [
        { path: 'terraform-provider-slack_v*.zip' },
        { path: 'terraform-provider-slack_v*_SHA256SUMS' },
        { path: 'terraform-provider-slack_v*_SHA256SUMS.sig' },
        { path: 'terraform-provider-slack_v*_manifest.json' },
      ]
    }],
    ['@semantic-release/git', {
      assets: ['go.mod', 'go.sum'],
      message: 'chore(release): ${nextRelease.version} [skip ci]\n\n${nextRelease.notes}',
    }],
  ]
};

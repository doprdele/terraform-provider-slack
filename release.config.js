module.exports = {
  branches: ['main'],
  plugins: [
    '@semantic-release/commit-analyzer',
    '@semantic-release/release-notes-generator',
    ['@semantic-release/exec', {
      prepareCmd: './build.sh ${nextRelease.version}',
    }],
    ['@semantic-release/github', {
      assets: [
        { path: 'terraform-provider-slack*.zip', label: 'Terraform Provider Binaries' },
      ]
    }],
    ['@semantic-release/git', {
      assets: ['go.mod', 'go.sum'],
      message: 'chore(release): ${nextRelease.version} [ci skip]\n\n${nextRelease.notes}'
    }]
  ]
};

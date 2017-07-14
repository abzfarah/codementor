// Copyright (c) 2018 Nomad Media, Inc. All Rights Reserved.


import assert from 'assert';

import * as Markdown from 'utils/markdown.jsx';

describe('Markdown.Imgs', function() {
    it('Inline mage', function(done) {
        assert.equal(
            Markdown.format('![Mattermost](/images/icon.png)').trim(),
            '<p><img src="/images/icon.png" alt="Mattermost" onload="window.markdownImageLoaded(this)" onerror="window.markdownImageLoaded(this)" class="markdown-inline-img"></p>'
        );

        done();
    });

    it('Image with hover text', function(done) {
        assert.equal(
            Markdown.format('![Mattermost](/images/icon.png "Mattermost Icon")').trim(),
            '<p><img src="/images/icon.png" alt="Mattermost" title="Mattermost Icon" onload="window.markdownImageLoaded(this)" onerror="window.markdownImageLoaded(this)" class="markdown-inline-img"></p>'
        );

        done();
    });

    it('Image with link', function(done) {
        assert.equal(
            Markdown.format('[![Mattermost](../../images/icon-76x76.png)](https://github.com/nomadsingles/platform)').trim(),
            '<p><a class="theme markdown__link" href="https://github.com/nomadsingles/platform" rel="noreferrer" target="_blank"><img src="../../images/icon-76x76.png" alt="Mattermost" onload="window.markdownImageLoaded(this)" onerror="window.markdownImageLoaded(this)" class="markdown-inline-img"></a></p>'
        );

        done();
    });

    it('Image with width and height', function(done) {
        assert.equal(
            Markdown.format('![Mattermost](../../images/icon-76x76.png =50x76 "Mattermost Icon")').trim(),
            '<p><img src="../../images/icon-76x76.png" alt="Mattermost" title="Mattermost Icon" width="50" height="76" onload="window.markdownImageLoaded(this)" onerror="window.markdownImageLoaded(this)" class="markdown-inline-img"></p>'
        );

        done();
    });

    it('Image with width', function(done) {
        assert.equal(
            Markdown.format('![Mattermost](../../images/icon-76x76.png =50 "Mattermost Icon")').trim(),
            '<p><img src="../../images/icon-76x76.png" alt="Mattermost" title="Mattermost Icon" width="50" onload="window.markdownImageLoaded(this)" onerror="window.markdownImageLoaded(this)" class="markdown-inline-img"></p>'
        );

        done();
    });
});

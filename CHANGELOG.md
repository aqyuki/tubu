# Changelog

## [1.1.0](https://github.com/aqyuki/tubu/compare/v1.0.1...v1.1.0) (2024-08-17)


### Features

* **bot:** change to allow recovery from the panic occurred ([#41](https://github.com/aqyuki/tubu/issues/41)) ([1a62701](https://github.com/aqyuki/tubu/commit/1a62701b2d73fcc18d37c62ede1b6046703c4be5))
* **command:** add `channel` command ([#37](https://github.com/aqyuki/tubu/issues/37)) ([263bfc3](https://github.com/aqyuki/tubu/commit/263bfc308c5ca2dd22035aa530814408e5b573df))
* **command:** add `guild` command ([#46](https://github.com/aqyuki/tubu/issues/46)) ([5dd8a3b](https://github.com/aqyuki/tubu/commit/5dd8a3be4488fad2cdb136130908f4ec4518fd32))
* **command:** add SendDM command ([#50](https://github.com/aqyuki/tubu/issues/50)) ([b92ff09](https://github.com/aqyuki/tubu/commit/b92ff09072cbc27600f4fe4aac5597fa5c17deb1))
* **command:** change maximum number of dice and maximum number of faces ([#48](https://github.com/aqyuki/tubu/issues/48)) ([46e0e4d](https://github.com/aqyuki/tubu/commit/46e0e4dfa25db21e71394b3bd1605fcfa0b98474))
* **internal:** change to allow register command multiple ([#53](https://github.com/aqyuki/tubu/issues/53)) ([90b35f6](https://github.com/aqyuki/tubu/commit/90b35f6f735d34ddbc4d6b7dc12d33747405da3a))


### Bug Fixes

* **docker:** fix the style of the embedded version ([#52](https://github.com/aqyuki/tubu/issues/52)) ([289767e](https://github.com/aqyuki/tubu/commit/289767ebf61f4f2fcb9d5e1a7d41853491df6b27))

## [1.0.1](https://github.com/aqyuki/tubu/compare/v1.0.0...v1.0.1) (2024-08-13)


### Bug Fixes

* fixed an issue where metadata could not be embedded in Docker images. ([#30](https://github.com/aqyuki/tubu/issues/30)) ([32a6838](https://github.com/aqyuki/tubu/commit/32a6838d9f83dcbf20ee6b99d98db95188012b16))

## 1.0.0 (2024-08-13)


### ⚠ BREAKING CHANGES

* **handler:** change path ([#16](https://github.com/aqyuki/tubu/issues/16))

### Features

* add a message expand handler ([#25](https://github.com/aqyuki/tubu/issues/25)) ([873500f](https://github.com/aqyuki/tubu/commit/873500faa6bb6a185b14369e0331a69b661ed64f))
* **bot:** add a Ready handler ([#11](https://github.com/aqyuki/tubu/issues/11)) ([363236d](https://github.com/aqyuki/tubu/commit/363236d3d38c5dff43c5b056dc03f61b2dd0fddd))
* **bot:** add entry point ([#12](https://github.com/aqyuki/tubu/issues/12)) ([3970daa](https://github.com/aqyuki/tubu/commit/3970daaac1b8d9bd0bb426e7baf7c8200ee1ba67))
* change to load configuration from environment variable ([#18](https://github.com/aqyuki/tubu/issues/18)) ([e7a6012](https://github.com/aqyuki/tubu/commit/e7a6012082051b561000cbde7a3a302d953178af))
* **command:** add `version` command ([#17](https://github.com/aqyuki/tubu/issues/17)) ([f417ef2](https://github.com/aqyuki/tubu/commit/f417ef252968ef9190a91ada947625c89a6e369b))
* **command:** add a `dice` command ([#27](https://github.com/aqyuki/tubu/issues/27)) ([97a3ba2](https://github.com/aqyuki/tubu/commit/97a3ba258b8ded5bf5ffad0e9be4b0f09d8b443c))
* **discord:** add a client to create a discord bot ([#8](https://github.com/aqyuki/tubu/issues/8)) ([e8f40a2](https://github.com/aqyuki/tubu/commit/e8f40a242afeb293f87edd13db65c9ac20d1425d))
* **discord:** register command router ([#10](https://github.com/aqyuki/tubu/issues/10)) ([904d4f4](https://github.com/aqyuki/tubu/commit/904d4f4e7e0b2b454dd5f89391cf871fef4e1d99))
* **docker:** add Dockerfile ([#20](https://github.com/aqyuki/tubu/issues/20)) ([5490392](https://github.com/aqyuki/tubu/commit/549039202b2b7c8ddbdbb6aa858264946cb64202))
* **logging:** add logging utilities ([#3](https://github.com/aqyuki/tubu/issues/3)) ([4a0ec3c](https://github.com/aqyuki/tubu/commit/4a0ec3cf273ad9eb1eab941b4197fbadef826214))
* **metadata:** change metadata package ([#14](https://github.com/aqyuki/tubu/issues/14)) ([73d1753](https://github.com/aqyuki/tubu/commit/73d1753b8f9ef8bba47a9866f03db5ea0ecd08bd))


### Miscellaneous Chores

* **handler:** change path ([#16](https://github.com/aqyuki/tubu/issues/16)) ([2616f1f](https://github.com/aqyuki/tubu/commit/2616f1f32744a2b00a15221f900d2d04df0bf6cb))

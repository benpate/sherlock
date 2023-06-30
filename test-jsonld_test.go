package sherlock

import (
	"bytes"
	"testing"

	"github.com/benpate/rosetta/mapof"
	"github.com/davecgh/go-spew/spew"
	"github.com/stretchr/testify/require"
)

func TestJSONLD(t *testing.T) {

	html := `<!DOCTYPE html>
<html lang='en'>
<head>
<meta charset='utf-8'>
<meta content='width=device-width, initial-scale=1' name='viewport'>
<link href='https://files.mastodon.social' rel='dns-prefetch'>
<link href='/favicon.ico' rel='icon' type='image/x-icon'>
<link href='/packs/media/icons/favicon-16x16-c58fdef40ced38d582d5b8eed9d15c5a.png' rel='icon' sizes='16x16' type='image/png'>
<link href='/packs/media/icons/favicon-32x32-249409a6d9f300112c51af514d863112.png' rel='icon' sizes='32x32' type='image/png'>
<link href='/packs/media/icons/favicon-48x48-c1197e9664ee6476d2715a1c4293bf61.png' rel='icon' sizes='48x48' type='image/png'>
<link href='/packs/media/icons/apple-touch-icon-57x57-c9dca808280860c51d0357f6a3350f4d.png' rel='apple-touch-icon' sizes='57x57'>
<link href='/packs/media/icons/apple-touch-icon-60x60-bb247db195d9ee9d8c687b2b048508d9.png' rel='apple-touch-icon' sizes='60x60'>
<link href='/packs/media/icons/apple-touch-icon-72x72-799d90b81f5b28cea7355a0c0b356381.png' rel='apple-touch-icon' sizes='72x72'>
<link href='/packs/media/icons/apple-touch-icon-76x76-015d73d770413d121873193153ae4ab5.png' rel='apple-touch-icon' sizes='76x76'>
<link href='/packs/media/icons/apple-touch-icon-114x114-211e68807b8d15707666a8d326d338b1.png' rel='apple-touch-icon' sizes='114x114'>
<link href='/packs/media/icons/apple-touch-icon-120x120-4c4e647d758bf1b2f47a53e2165a01d2.png' rel='apple-touch-icon' sizes='120x120'>
<link href='/packs/media/icons/apple-touch-icon-144x144-ff3110f7772743bdd0c1c47fb7b2d4e0.png' rel='apple-touch-icon' sizes='144x144'>
<link href='/packs/media/icons/apple-touch-icon-152x152-b12cbb1baaf4a6111d0efd391fd829c9.png' rel='apple-touch-icon' sizes='152x152'>
<link href='/packs/media/icons/apple-touch-icon-167x167-6f21a68f6a05a8b5cf25b1729e553728.png' rel='apple-touch-icon' sizes='167x167'>
<link href='/packs/media/icons/apple-touch-icon-180x180-a75559a0af48064c1b7c71b81f3bf7c6.png' rel='apple-touch-icon' sizes='180x180'>
<link href='/packs/media/icons/apple-touch-icon-1024x1024-db6849588b44f525363c37b65ef0ac66.png' rel='apple-touch-icon' sizes='1024x1024'>
<link color='#6364FF' href='/packs/media/images/logo-symbol-icon-de9e68dc49b19eb5cd142d3316f9e95e.svg' rel='mask-icon'>
<link href='/manifest' rel='manifest'>
<meta content='#191b22' name='theme-color'>
<meta content='yes' name='apple-mobile-web-app-capable'>
<title>Ben Pate 🤘: &quot;So.. Twitter has been up and down for me this aft…&quot; - Mastodon</title>
<link rel="stylesheet" media="all" crossorigin="anonymous" href="/packs/css/common-a729b6b0.css" integrity="sha256-KDzFV9ckqn2zELncHPapWY+nN4HgaUP+jxBFg4sinFA=" />
<link rel="stylesheet" media="all" crossorigin="anonymous" href="/packs/css/default-db340b81.chunk.css" integrity="sha256-kz+i8IxlKdZtYp6OypigB9sPxapI3iHpSRGqYl+t8VY=" />
<script src="/packs/js/common-2605ec533371c6e88d51.js" crossorigin="anonymous" integrity="sha256-y9pATNAYdqblI4mAcoAYkQDRw8MxyNx3l+gifzhQas0="></script>
<script src="/packs/js/locale_en-afe3f6fbe49cbd149fe6.chunk.js" crossorigin="anonymous" integrity="sha256-58bCBNAWxQi33l9kFF8Qx68UylHJlBc0wRTIk0x3nSw="></script>

<meta content='N+1+Y0yCmdjJsWvGX8ASNQ==' name='style-nonce'>
<link rel="stylesheet" media="all" href="/inert.css" id="inert-style" />
<link rel="stylesheet" media="all" href="https://mastodon.social/custom.css" />
<link href='https://mastodon.social/api/oembed?format=json&amp;url=https%3A%2F%2Fmastodon.social%2F%40benpate%2F109684236270111476' rel='alternate' type='application/json+oembed'>
<link href='https://mastodon.social/users/benpate/statuses/109684236270111476' rel='alternate' type='application/activity+json'>
<meta content="Mastodon" property="og:site_name">
<meta content="article" property="og:type">
<meta content="Ben Pate 🤘 (@benpate@mastodon.social)" property="og:title">
<meta content="https://mastodon.social/@benpate/109684236270111476" property="og:url">
<meta content="2023-01-13T22:23:44Z" property="og:published_time">
<meta content="benpate@mastodon.social" property="profile:username">
<meta content='So.. Twitter has been up and down for me this afternoon.  I don&#39;t think this is an &quot;it&#39;s happening&quot; moment, but it&#39;s definitely not a good look for corporate social media.  

I can&#39;t wait to have my own ActivityPub app online on the Fediverse!' name='description'>
<meta content="So.. Twitter has been up and down for me this afternoon.  I don&#39;t think this is an &quot;it&#39;s happening&quot; moment, but it&#39;s definitely not a good look for corporate social media.  

I can&#39;t wait to have my own ActivityPub app online on the Fediverse!" property="og:description">

<meta content="summary" property="twitter:card">

<meta content='BCk-QqERU0q-CfYZjcuB6lnyyOYfJ2AifKqfeGIm7Z-HiTU5T9eTG5GxVA0_OH5mMlI4UkkDTpaZwozy0TzdZ2M=' name='applicationServerKey'>
<script id="initial-state" type="application/json">{"meta":{"streaming_api_base_url":"wss://streaming.mastodon.social","access_token":null,"locale":"en","domain":"mastodon.social","title":"Mastodon","admin":null,"search_enabled":true,"repository":"mastodon/mastodon","source_url":"https://github.com/mastodon/mastodon","version":"4.1.2+nightly-20230523","limited_federation_mode":false,"mascot":"https://files.mastodon.social/site_uploads/files/000/000/002/original/080cba45af6a6356.svg","profile_directory":true,"trends":true,"registrations_open":true,"timeline_preview":true,"activity_api_enabled":true,"single_user_mode":false,"trends_as_landing_page":true,"status_page_url":"https://status.mastodon.social","auto_play_gif":null,"display_media":null,"reduce_motion":null,"use_blurhash":null,"crop_images":null},"compose":{"text":""},"accounts":{},"media_attachments":{"accept_content_types":[".jpg",".jpeg",".png",".gif",".webp",".heic",".heif",".avif",".webm",".mp4",".m4v",".mov",".ogg",".oga",".mp3",".wav",".flac",".opus",".aac",".m4a",".3gp",".wma","image/jpeg","image/png","image/gif","image/heic","image/heif","image/webp","image/avif","video/webm","video/mp4","video/quicktime","video/ogg","audio/wave","audio/wav","audio/x-wav","audio/x-pn-wave","audio/vnd.wave","audio/ogg","audio/vorbis","audio/mpeg","audio/mp3","audio/webm","audio/flac","audio/aac","audio/m4a","audio/x-m4a","audio/mp4","audio/3gpp","video/x-ms-asf"]},"settings":{},"languages":[["aa","Afar","Afaraf"],["ab","Abkhaz","аҧсуа бызшәа"],["ae","Avestan","avesta"],["af","Afrikaans","Afrikaans"],["ak","Akan","Akan"],["am","Amharic","አማርኛ"],["an","Aragonese","aragonés"],["ar","Arabic","اللغة العربية"],["as","Assamese","অসমীয়া"],["av","Avaric","авар мацӀ"],["ay","Aymara","aymar aru"],["az","Azerbaijani","azərbaycan dili"],["ba","Bashkir","башҡорт теле"],["be","Belarusian","беларуская мова"],["bg","Bulgarian","български език"],["bh","Bihari","भोजपुरी"],["bi","Bislama","Bislama"],["bm","Bambara","bamanankan"],["bn","Bengali","বাংলা"],["bo","Tibetan","བོད་ཡིག"],["br","Breton","brezhoneg"],["bs","Bosnian","bosanski jezik"],["ca","Catalan","Català"],["ce","Chechen","нохчийн мотт"],["ch","Chamorro","Chamoru"],["co","Corsican","corsu"],["cr","Cree","ᓀᐦᐃᔭᐍᐏᐣ"],["cs","Czech","čeština"],["cu","Old Church Slavonic","ѩзыкъ словѣньскъ"],["cv","Chuvash","чӑваш чӗлхи"],["cy","Welsh","Cymraeg"],["da","Danish","dansk"],["de","German","Deutsch"],["dv","Divehi","Dhivehi"],["dz","Dzongkha","རྫོང་ཁ"],["ee","Ewe","Eʋegbe"],["el","Greek","Ελληνικά"],["en","English","English"],["eo","Esperanto","Esperanto"],["es","Spanish","Español"],["et","Estonian","eesti"],["eu","Basque","euskara"],["fa","Persian","فارسی"],["ff","Fula","Fulfulde"],["fi","Finnish","suomi"],["fj","Fijian","Vakaviti"],["fo","Faroese","føroyskt"],["fr","French","Français"],["fy","Western Frisian","Frysk"],["ga","Irish","Gaeilge"],["gd","Scottish Gaelic","Gàidhlig"],["gl","Galician","galego"],["gu","Gujarati","ગુજરાતી"],["gv","Manx","Gaelg"],["ha","Hausa","هَوُسَ"],["he","Hebrew","עברית"],["hi","Hindi","हिन्दी"],["ho","Hiri Motu","Hiri Motu"],["hr","Croatian","Hrvatski"],["ht","Haitian","Kreyòl ayisyen"],["hu","Hungarian","magyar"],["hy","Armenian","Հայերեն"],["hz","Herero","Otjiherero"],["ia","Interlingua","Interlingua"],["id","Indonesian","Bahasa Indonesia"],["ie","Interlingue","Interlingue"],["ig","Igbo","Asụsụ Igbo"],["ii","Nuosu","ꆈꌠ꒿ Nuosuhxop"],["ik","Inupiaq","Iñupiaq"],["io","Ido","Ido"],["is","Icelandic","Íslenska"],["it","Italian","Italiano"],["iu","Inuktitut","ᐃᓄᒃᑎᑐᑦ"],["ja","Japanese","日本語"],["jv","Javanese","basa Jawa"],["ka","Georgian","ქართული"],["kg","Kongo","Kikongo"],["ki","Kikuyu","Gĩkũyũ"],["kj","Kwanyama","Kuanyama"],["kk","Kazakh","қазақ тілі"],["kl","Kalaallisut","kalaallisut"],["km","Khmer","ខេមរភាសា"],["kn","Kannada","ಕನ್ನಡ"],["ko","Korean","한국어"],["kr","Kanuri","Kanuri"],["ks","Kashmiri","कश्मीरी"],["ku","Kurmanji (Kurdish)","Kurmancî"],["kv","Komi","коми кыв"],["kw","Cornish","Kernewek"],["ky","Kyrgyz","Кыргызча"],["la","Latin","latine"],["lb","Luxembourgish","Lëtzebuergesch"],["lg","Ganda","Luganda"],["li","Limburgish","Limburgs"],["ln","Lingala","Lingála"],["lo","Lao","ລາວ"],["lt","Lithuanian","lietuvių kalba"],["lu","Luba-Katanga","Tshiluba"],["lv","Latvian","latviešu valoda"],["mg","Malagasy","fiteny malagasy"],["mh","Marshallese","Kajin M̧ajeļ"],["mi","Māori","te reo Māori"],["mk","Macedonian","македонски јазик"],["ml","Malayalam","മലയാളം"],["mn","Mongolian","Монгол хэл"],["mr","Marathi","मराठी"],["ms","Malay","Bahasa Melayu"],["mt","Maltese","Malti"],["my","Burmese","ဗမာစာ"],["na","Nauru","Ekakairũ Naoero"],["nb","Norwegian Bokmål","Norsk bokmål"],["nd","Northern Ndebele","isiNdebele"],["ne","Nepali","नेपाली"],["ng","Ndonga","Owambo"],["nl","Dutch","Nederlands"],["nn","Norwegian Nynorsk","Norsk Nynorsk"],["no","Norwegian","Norsk"],["nr","Southern Ndebele","isiNdebele"],["nv","Navajo","Diné bizaad"],["ny","Chichewa","chiCheŵa"],["oc","Occitan","occitan"],["oj","Ojibwe","ᐊᓂᔑᓈᐯᒧᐎᓐ"],["om","Oromo","Afaan Oromoo"],["or","Oriya","ଓଡ଼ିଆ"],["os","Ossetian","ирон æвзаг"],["pa","Panjabi","ਪੰਜਾਬੀ"],["pi","Pāli","पाऴि"],["pl","Polish","Polski"],["ps","Pashto","پښتو"],["pt","Portuguese","Português"],["qu","Quechua","Runa Simi"],["rm","Romansh","rumantsch grischun"],["rn","Kirundi","Ikirundi"],["ro","Romanian","Română"],["ru","Russian","Русский"],["rw","Kinyarwanda","Ikinyarwanda"],["sa","Sanskrit","संस्कृतम्"],["sc","Sardinian","sardu"],["sd","Sindhi","सिन्धी"],["se","Northern Sami","Davvisámegiella"],["sg","Sango","yângâ tî sängö"],["si","Sinhala","සිංහල"],["sk","Slovak","slovenčina"],["sl","Slovenian","slovenščina"],["sn","Shona","chiShona"],["so","Somali","Soomaaliga"],["sq","Albanian","Shqip"],["sr","Serbian","српски језик"],["ss","Swati","SiSwati"],["st","Southern Sotho","Sesotho"],["su","Sundanese","Basa Sunda"],["sv","Swedish","Svenska"],["sw","Swahili","Kiswahili"],["ta","Tamil","தமிழ்"],["te","Telugu","తెలుగు"],["tg","Tajik","тоҷикӣ"],["th","Thai","ไทย"],["ti","Tigrinya","ትግርኛ"],["tk","Turkmen","Türkmen"],["tl","Tagalog","Wikang Tagalog"],["tn","Tswana","Setswana"],["to","Tonga","faka Tonga"],["tr","Turkish","Türkçe"],["ts","Tsonga","Xitsonga"],["tt","Tatar","татар теле"],["tw","Twi","Twi"],["ty","Tahitian","Reo Tahiti"],["ug","Uyghur","ئۇيغۇرچە‎"],["uk","Ukrainian","Українська"],["ur","Urdu","اردو"],["uz","Uzbek","Ўзбек"],["ve","Venda","Tshivenḓa"],["vi","Vietnamese","Tiếng Việt"],["vo","Volapük","Volapük"],["wa","Walloon","walon"],["wo","Wolof","Wollof"],["xh","Xhosa","isiXhosa"],["yi","Yiddish","ייִדיש"],["yo","Yoruba","Yorùbá"],["za","Zhuang","Saɯ cueŋƅ"],["zh","Chinese","中文"],["zu","Zulu","isiZulu"],["ast","Asturian","Asturianu"],["ckb","Sorani (Kurdish)","سۆرانی"],["cnr","Montenegrin","crnogorski"],["jbo","Lojban","la .lojban."],["kab","Kabyle","Taqbaylit"],["kmr","Kurmanji (Kurdish)","Kurmancî"],["ldn","Láadan","Láadan"],["lfn","Lingua Franca Nova","lingua franca nova"],["sco","Scots","Scots"],["sma","Southern Sami","Åarjelsaemien Gïele"],["smj","Lule Sami","Julevsámegiella"],["szl","Silesian","ślůnsko godka"],["tok","Toki Pona","toki pona"],["zba","Balaibalan","باليبلن"],["zgh","Standard Moroccan Tamazight","ⵜⴰⵎⴰⵣⵉⵖⵜ"]],"push_subscription":null,"role":null}</script>
<script src="/packs/js/application-fdb624b6ff7af0b7b1ba.chunk.js" crossorigin="anonymous" integrity="sha256-OIjyL0rTQkC2zzphOD5YIDylvxK6SDzQbUNno+zdDRk="></script>

</head>
<body class='app-body theme-default no-reduce-motion'>
<div class='notranslate app-holder' data-props='{&quot;locale&quot;:&quot;en&quot;}' id='mastodon'>
<noscript>
<img alt="Mastodon" src="/packs/media/images/logo-d4b5dc90fd3e117d141ae7053b157f58.svg" />
<div>
To use the Mastodon web application, please enable JavaScript. Alternatively, try one of the <a href="https://joinmastodon.org/apps">native apps</a> for Mastodon for your platform.
</div>
</noscript>
</div>


<div aria-hidden class='logo-resources' inert tabindex='-1'>
<svg xmlns="http://www.w3.org/2000/svg" xmlns:xlink="http://www.w3.org/1999/xlink" width="79" height="79" viewBox="0 0 79 75"><symbol id="logo-symbol-icon"><path d="M74.7135 16.6043C73.6199 8.54587 66.5351 2.19527 58.1366 0.964691C56.7196 0.756754 51.351 0 38.9148 0H38.822C26.3824 0 23.7135 0.756754 22.2966 0.964691C14.1319 2.16118 6.67571 7.86752 4.86669 16.0214C3.99657 20.0369 3.90371 24.4888 4.06535 28.5726C4.29578 34.4289 4.34049 40.275 4.877 46.1075C5.24791 49.9817 5.89495 53.8251 6.81328 57.6088C8.53288 64.5968 15.4938 70.4122 22.3138 72.7848C29.6155 75.259 37.468 75.6697 44.9919 73.971C45.8196 73.7801 46.6381 73.5586 47.4475 73.3063C49.2737 72.7302 51.4164 72.086 52.9915 70.9542C53.0131 70.9384 53.0308 70.9178 53.0433 70.8942C53.0558 70.8706 53.0628 70.8445 53.0637 70.8179V65.1661C53.0634 65.1412 53.0574 65.1167 53.0462 65.0944C53.035 65.0721 53.0189 65.0525 52.9992 65.0371C52.9794 65.0218 52.9564 65.011 52.9318 65.0056C52.9073 65.0002 52.8819 65.0003 52.8574 65.0059C48.0369 66.1472 43.0971 66.7193 38.141 66.7103C29.6118 66.7103 27.3178 62.6981 26.6609 61.0278C26.1329 59.5842 25.7976 58.0784 25.6636 56.5486C25.6622 56.5229 25.667 56.4973 25.6775 56.4738C25.688 56.4502 25.7039 56.4295 25.724 56.4132C25.7441 56.397 25.7678 56.3856 25.7931 56.3801C25.8185 56.3746 25.8448 56.3751 25.8699 56.3816C30.6101 57.5151 35.4693 58.0873 40.3455 58.086C41.5183 58.086 42.6876 58.086 43.8604 58.0553C48.7647 57.919 53.9339 57.6701 58.7591 56.7361C58.8794 56.7123 58.9998 56.6918 59.103 56.6611C66.7139 55.2124 73.9569 50.665 74.6929 39.1501C74.7204 38.6967 74.7892 34.4016 74.7892 33.9312C74.7926 32.3325 75.3085 22.5901 74.7135 16.6043ZM62.9996 45.3371H54.9966V25.9069C54.9966 21.8163 53.277 19.7302 49.7793 19.7302C45.9343 19.7302 44.0083 22.1981 44.0083 27.0727V37.7082H36.0534V27.0727C36.0534 22.1981 34.124 19.7302 30.279 19.7302C26.8019 19.7302 25.0651 21.8163 25.0617 25.9069V45.3371H17.0656V25.3172C17.0656 21.2266 18.1191 17.9769 20.2262 15.568C22.3998 13.1648 25.2509 11.9308 28.7898 11.9308C32.8859 11.9308 35.9812 13.492 38.0447 16.6111L40.036 19.9245L42.0308 16.6111C44.0943 13.492 47.1896 11.9308 51.2788 11.9308C54.8143 11.9308 57.6654 13.1648 59.8459 15.568C61.9529 17.9746 63.0065 21.2243 63.0065 25.3172L62.9996 45.3371Z" fill="currentColor"/></symbol><use xlink:href="#logo-symbol-icon"/></svg>


<svg width="261" height="66" viewBox="0 0 261 66" fill="none" xmlns="http://www.w3.org/2000/svg" xmlns:xlink="http://www.w3.org/1999/xlink">
<symbol id="logo-symbol-wordmark"><path d="M60.7539 14.4034C59.8143 7.41942 53.7273 1.91557 46.5117 0.849066C45.2943 0.668854 40.6819 0.0130005 29.9973 0.0130005H29.9175C19.2299 0.0130005 16.937 0.668854 15.7196 0.849066C8.70488 1.88602 2.29885 6.83152 0.744617 13.8982C-0.00294988 17.3784 -0.0827298 21.2367 0.0561464 24.7759C0.254119 29.8514 0.292531 34.918 0.753482 39.9728C1.07215 43.3305 1.62806 46.6614 2.41704 49.9406C3.89445 55.9969 9.87499 61.0369 15.7344 63.0931C22.0077 65.2374 28.7542 65.5934 35.2184 64.1212C35.9295 63.9558 36.6318 63.7638 37.3252 63.5451C38.8971 63.0459 40.738 62.4875 42.0913 61.5067C42.1099 61.4929 42.1251 61.4751 42.1358 61.4547C42.1466 61.4342 42.1526 61.4116 42.1534 61.3885V56.4903C42.153 56.4687 42.1479 56.4475 42.1383 56.4281C42.1287 56.4088 42.1149 56.3918 42.0979 56.3785C42.0809 56.3652 42.0611 56.3559 42.04 56.3512C42.019 56.3465 41.9971 56.3466 41.9761 56.3514C37.8345 57.3406 33.5905 57.8364 29.3324 57.8286C22.0045 57.8286 20.0336 54.3514 19.4693 52.9038C19.0156 51.6527 18.7275 50.3476 18.6124 49.0218C18.6112 48.9996 18.6153 48.9773 18.6243 48.9569C18.6333 48.9366 18.647 48.9186 18.6643 48.9045C18.6816 48.8904 18.7019 48.8805 18.7237 48.8758C18.7455 48.871 18.7681 48.8715 18.7897 48.8771C22.8622 49.8595 27.037 50.3553 31.2265 50.3542C32.234 50.3542 33.2387 50.3542 34.2463 50.3276C38.4598 50.2094 42.9009 49.9938 47.0465 49.1843C47.1499 49.1636 47.2534 49.1459 47.342 49.1193C53.881 47.8637 60.1038 43.9227 60.7362 33.9431C60.7598 33.5502 60.8189 29.8278 60.8189 29.4201C60.8218 28.0345 61.2651 19.5911 60.7539 14.4034Z" fill="url(#paint0_linear_89_11)"/>
<path d="M12.3442 18.3034C12.3442 16.2668 13.9777 14.6194 15.997 14.6194C18.0163 14.6194 19.6497 16.2668 19.6497 18.3034C19.6497 20.34 18.0163 21.9874 15.997 21.9874C13.9777 21.9874 12.3442 20.34 12.3442 18.3034Z" fill="currentColor"/>
<path d="M66.1484 21.4685V38.3839H59.4988V21.9744C59.4988 18.5109 58.0583 16.7597 55.1643 16.7597C51.9746 16.7597 50.3668 18.8482 50.3668 22.9603V31.9499H43.7687V22.9603C43.7687 18.8352 42.1738 16.7597 38.9712 16.7597C36.0901 16.7597 34.6367 18.5109 34.6367 21.9744V38.3839H28V21.4685C28 18.018 28.8746 15.268 30.6238 13.2314C32.4374 11.1948 34.8039 10.157 37.7365 10.157C41.132 10.157 43.7172 11.4802 45.415 14.1135L47.0742 16.9154L48.7334 14.1135C50.4311 11.4802 53.0035 10.157 56.4119 10.157C59.3444 10.157 61.711 11.1948 63.5246 13.2314C65.2738 15.268 66.1484 18.005 66.1484 21.4685ZM89.0297 29.8743C90.4059 28.4085 91.0619 26.5795 91.0619 24.3613C91.0619 22.1431 90.4059 20.3011 89.0297 18.9001C87.7049 17.4343 86.0329 16.7338 84.0007 16.7338C81.9685 16.7338 80.2965 17.4343 78.9717 18.9001C77.6469 20.3011 76.991 22.1431 76.991 24.3613C76.991 26.5795 77.6469 28.4215 78.9717 29.8743C80.2965 31.2753 81.9685 31.9888 84.0007 31.9888C86.0329 31.9888 87.7049 31.2883 89.0297 29.8743ZM91.0619 10.8316H97.6086V37.891H91.0619V34.6999C89.0811 37.3462 86.3416 38.6563 82.7788 38.6563C79.2161 38.6563 76.4765 37.3073 74.0456 34.5442C71.6533 31.7812 70.4443 28.3696 70.4443 24.3743C70.4443 20.3789 71.6661 17.0192 74.0456 14.2561C76.4893 11.4931 79.3833 10.0922 82.7788 10.0922C86.1744 10.0922 89.0811 11.3894 91.0619 14.0356V10.8445V10.8316ZM119.654 23.8683C121.583 25.3342 122.548 27.3837 122.496 29.9781C122.496 32.7411 121.532 34.9075 119.551 36.4122C117.57 37.878 115.178 38.6304 112.284 38.6304C107.049 38.6304 103.499 36.4641 101.621 32.1963L107.306 28.7847C108.065 31.1067 109.737 32.3001 112.284 32.3001C114.625 32.3001 115.782 31.5477 115.782 29.9781C115.782 28.8366 114.265 27.8118 111.165 27.0075C109.995 26.6833 109.03 26.359 108.271 26.0865C107.204 25.6585 106.29 25.1655 105.532 24.5688C103.654 23.103 102.689 21.1572 102.689 18.6666C102.689 16.0203 103.602 13.9059 105.429 12.3882C107.306 10.8186 109.596 10.0662 112.335 10.0662C116.709 10.0662 119.898 11.9601 121.982 15.7998L116.4 19.0428C115.59 17.2008 114.213 16.2798 112.335 16.2798C110.355 16.2798 109.39 17.0321 109.39 18.498C109.39 19.6395 110.908 20.6643 114.008 21.4685C116.4 22.0134 118.278 22.8176 119.641 23.8554L119.654 23.8683ZM140.477 17.538H134.741V28.7977C134.741 30.1468 135.255 30.964 136.22 31.3402C136.927 31.6126 138.355 31.6645 140.49 31.5607V37.891C136.079 38.4358 132.876 37.9948 130.998 36.5419C129.12 35.1409 128.207 32.5336 128.207 28.8106V17.538H123.795V10.8316H128.207V5.37038L134.754 3.25595V10.8316H140.49V17.538H140.477ZM161.352 29.7187C162.677 28.3177 163.333 26.5276 163.333 24.3613C163.333 22.195 162.677 20.4178 161.352 19.0039C160.027 17.6029 158.407 16.8894 156.426 16.8894C154.445 16.8894 152.825 17.5899 151.5 19.0039C150.227 20.4697 149.571 22.2469 149.571 24.3613C149.571 26.4757 150.227 28.2529 151.5 29.7187C152.825 31.1196 154.445 31.8331 156.426 31.8331C158.407 31.8331 160.027 31.1326 161.352 29.7187ZM146.883 34.5313C144.297 31.7682 143.024 28.4215 143.024 24.3613C143.024 20.3011 144.297 17.0062 146.883 14.2432C149.468 11.4802 152.67 10.0792 156.426 10.0792C160.182 10.0792 163.384 11.4802 165.97 14.2432C168.555 17.0062 169.88 20.4178 169.88 24.3613C169.88 28.3047 168.555 31.7682 165.97 34.5313C163.384 37.2943 160.233 38.6434 156.426 38.6434C152.619 38.6434 149.468 37.2943 146.883 34.5313ZM191.771 29.8743C193.095 28.4085 193.751 26.5795 193.751 24.3613C193.751 22.1431 193.095 20.3011 191.771 18.9001C190.446 17.4343 188.774 16.7338 186.742 16.7338C184.709 16.7338 183.037 17.4343 181.661 18.9001C180.336 20.3011 179.68 22.1431 179.68 24.3613C179.68 26.5795 180.336 28.4215 181.661 29.8743C183.037 31.2753 184.761 31.9888 186.742 31.9888C188.722 31.9888 190.446 31.2883 191.771 29.8743ZM193.751 0H200.298V37.891H193.751V34.6999C191.822 37.3462 189.082 38.6563 185.52 38.6563C181.957 38.6563 179.179 37.3073 176.735 34.5442C174.343 31.7812 173.134 28.3696 173.134 24.3743C173.134 20.3789 174.356 17.0192 176.735 14.2561C179.166 11.4931 182.111 10.0922 185.52 10.0922C188.928 10.0922 191.822 11.3894 193.751 14.0356V0.0129719V0ZM223.308 29.7057C224.633 28.3047 225.289 26.5146 225.289 24.3483C225.289 22.182 224.633 20.4048 223.308 18.9909C221.983 17.5899 220.363 16.8765 218.382 16.8765C216.401 16.8765 214.78 17.577 213.456 18.9909C212.182 20.4567 211.526 22.2339 211.526 24.3483C211.526 26.4627 212.182 28.2399 213.456 29.7057C214.78 31.1067 216.401 31.8201 218.382 31.8201C220.363 31.8201 221.983 31.1196 223.308 29.7057ZM208.838 34.5183C206.253 31.7553 204.98 28.4085 204.98 24.3483C204.98 20.2881 206.253 16.9932 208.838 14.2302C211.424 11.4672 214.626 10.0662 218.382 10.0662C222.137 10.0662 225.34 11.4672 227.925 14.2302C230.511 16.9932 231.835 20.4048 231.835 24.3483C231.835 28.2918 230.511 31.7553 227.925 34.5183C225.34 37.2813 222.189 38.6304 218.382 38.6304C214.575 38.6304 211.424 37.2813 208.838 34.5183ZM260.17 21.261V37.878H253.623V22.1301C253.623 20.34 253.173 18.9909 252.247 17.9661C251.385 17.0451 250.164 16.5651 248.594 16.5651C244.89 16.5651 243.012 18.7833 243.012 23.2716V37.878H236.466V10.8316H243.012V13.867C244.581 11.3245 247.077 10.0792 250.575 10.0792C253.366 10.0792 255.656 11.0521 257.431 13.0498C259.257 15.0474 260.17 17.7586 260.17 21.274" fill="currentColor"/>
<defs>
<linearGradient id="paint0_linear_89_11" x1="30.5" y1="0.0130005" x2="30.5" y2="65.013" gradientUnits="userSpaceOnUse">
<stop stop-color="#6364FF"/>
<stop offset="1" stop-color="#563ACC"/>
</linearGradient>
</defs></symbol><use xlink:href="#logo-symbol-wordmark"/>
</svg>

</div>
</body>
</html>
`

	json := mapof.NewAny()
	err := Parse("https://indieweb.org", bytes.NewBufferString(html), json)

	require.Nil(t, err)
	spew.Dump(json)
}

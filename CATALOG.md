# GigMCP Manifest Catalog

This catalog mirrors [Composio's toolkit list](https://composio.dev/tools) as a curated, lint-enforced set of MCP-server manifests. Most entries are **toolspec-driven**: a `manifests/<name>/<version>.toolspec.yaml` maps each manifest tool to a real REST endpoint, served by the generic [toolpack](https://github.com/gigmcp/toolpack) engine. Every toolspec was researched against the service's official API docs and adversarially re-verified by an independent pass; `toolspec (unverified)` marks the few whose docs could not be independently confirmed. `adopts ...` entries build an established upstream Go MCP server instead. Image digests are placeholders (`sha256:0000...`) until `build-images` CI pins the real ones — placeholder digests are **not installable**.

| Name | Auth | Tier | Egress | Status |
|------|------|------|--------|--------|
| ably | api_key | sealed | `rest.ably.io` | toolspec |
| acculynx | api_key | sealed | `api.acculynx.com` | toolspec |
| activecampaign | api_key | sealed | `*.api-us1.com` | toolspec |
| affinity | api_key | sealed | `api.affinity.co` | toolspec |
| agencyzoom | api_key | sealed | `api.agencyzoom.com` | toolspec |
| ahrefs | api_key | sealed | `api.ahrefs.com` | toolspec |
| airtable | oauth2 | sealed | `api.airtable.com` | toolspec |
| alchemy | api_key | entrusted | `*.g.alchemy.com` | toolspec |
| altoviz | api_key | sealed | `api.altoviz.com` | toolspec (unverified) |
| amcards | api_key | sealed | `amcards.com` | toolspec |
| amplitude | basic | sealed | `amplitude.com`, `api2.amplitude.com` | toolspec |
| apaleo | oauth2 | sealed | `api.apaleo.com`, `identity.apaleo.com` | toolspec |
| apollo | api_key | sealed | `api.apollo.io` | toolspec |
| appdrag | api_key | sealed | `api.appdrag.com` | toolspec |
| asana | oauth2 | sealed | `app.asana.com` | toolspec |
| ashby | basic | sealed | `api.ashbyhq.com` | toolspec |
| attio | oauth2 | sealed | `api.attio.com` | toolspec |
| bamboohr | basic | sealed | `api.bamboohr.com` | toolspec |
| bannerbear | api_key | sealed | `api.bannerbear.com` | toolspec |
| baserow | api_key | sealed | `api.baserow.io` | toolspec |
| beeminder | api_key | entrusted | `www.beeminder.com` | toolspec |
| bitbucket | oauth2 | sealed | `api.bitbucket.org`, `bitbucket.org` | toolspec |
| bitwarden | oauth2 | sealed | `api.bitwarden.com`, `identity.bitwarden.com` | toolspec |
| blackbaud | oauth2 | sealed | `api.sky.blackbaud.com`, `oauth2.sky.blackbaud.com` | toolspec |
| blackboard | oauth2 | sealed | `*.blackboard.com` | toolspec |
| boldsign | oauth2 | sealed | `api.boldsign.com` | toolspec |
| bolna | api_key | sealed | `api.bolna.ai` | toolspec |
| borneo | api_key | sealed | `api.borneo.io`, `*.borneo.io` | planned (no public API found) |
| botbaba | api_key | sealed | `botbaba.io`, `*.botbaba.io` | planned (no public API found) |
| box | oauth2 | sealed | `api.box.com`, `upload.box.com` | toolspec |
| brandfetch | api_key | sealed | `api.brandfetch.io`, `graphql.brandfetch.io` | toolspec |
| breezy-hr | api_key | sealed | `api.breezy.hr` | toolspec |
| brevo | api_key | sealed | `api.brevo.com` | toolspec |
| brex | oauth2 | sealed | `api.brex.com` | toolspec |
| browseai | api_key | sealed | `api.browse.ai` | toolspec |
| browserbase | api_key | sealed | `api.browserbase.com` | toolspec |
| bubble | api_key | sealed | `bubble.io`, `*.bubbleapps.io` | toolspec |
| cal | api_key | sealed | `api.cal.com` | toolspec |
| calendly | oauth2 | sealed | `api.calendly.com`, `auth.calendly.com` | toolspec |
| canva | oauth2 | sealed | `api.canva.com` | toolspec |
| canvas | api_key | sealed | `canvas.instructure.com`, `*.instructure.com` | toolspec |
| chatwork | api_key | sealed | `api.chatwork.com` | toolspec |
| chmeetings | api_key | sealed | `api.chmeetings.com`, `chmeetings.com` | planned (no public API found) |
| clickup | oauth2 | sealed | `api.clickup.com` | toolspec |
| close | basic | sealed | `api.close.com` | toolspec |
| cloudflare | api_key | sealed | `api.cloudflare.com` | toolspec |
| coda | api_key | sealed | `coda.io` | toolspec |
| coinbase | api_key | sealed | `api.coinbase.com`, `api.exchange.coinbase.com` | toolspec |
| coinmarketcal | api_key | sealed | `developers.coinmarketcal.com` | toolspec |
| confluence | oauth2 | sealed | `api.atlassian.com`, `*.atlassian.net` | toolspec |
| contentful | oauth2 | sealed | `api.contentful.com` | toolspec |
| crustdata | api_key | sealed | `api.crustdata.com` | toolspec (unverified) |
| d2lbrightspace | oauth2 | sealed | `auth.brightspace.com`, `*.brightspace.com` | toolspec |
| dailybot | api_key | sealed | `api.dailybot.com` | toolspec |
| datadog | api_key | sealed | `api.datadoghq.com`, `*.datadoghq.com` | toolspec |
| datagma | api_key | entrusted | `gateway.datagma.net` | toolspec (unverified) |
| datarobot | api_key | sealed | `app.datarobot.com` | toolspec |
| demio | api_key | sealed | `my.demio.com` | toolspec (unverified) |
| dialpad | oauth2 | sealed | `dialpad.com` | toolspec |
| digicert | api_key | sealed | `www.digicert.com`, `api.digicert.com` | toolspec |
| discord | oauth2 | sealed | `discord.com` | toolspec |
| discordbot | api_key | sealed | `discord.com`, `cdn.discordapp.com` | toolspec |
| docmosis | api_key | sealed | `*.dws4.docmosis.com`, `*.dws3.docmosis.com` | toolspec |
| docusign | oauth2 | sealed | `account.docusign.com`, `*.docusign.net` | toolspec |
| dropbox | oauth2 | sealed | `api.dropboxapi.com`, `content.dropboxapi.com` | toolspec |
| dropbox-sign | oauth2 | sealed | `api.hellosign.com` | toolspec |
| dynamics365 | oauth2 | sealed | `*.dynamics.com`, `login.microsoftonline.com` | toolspec |
| echtpost | api_key | sealed | `api.echtpost.de` | toolspec |
| elevenlabs | api_key | sealed | `api.elevenlabs.io` | toolspec |
| eventbrite | oauth2 | sealed | `www.eventbriteapi.com` | toolspec |
| exa | api_key | sealed | `api.exa.ai` | toolspec |
| exist | oauth2 | sealed | `exist.io` | toolspec |
| facebook | oauth2 | sealed | `graph.facebook.com` | toolspec |
| figma | oauth2 | sealed | `api.figma.com` | toolspec |
| finage | api_key | entrusted | `api.finage.co.uk` | toolspec |
| firecrawl | api_key | sealed | `api.firecrawl.dev` | toolspec |
| fireflies | api_key | sealed | `api.fireflies.ai` | toolspec |
| flutterwave | api_key | sealed | `api.flutterwave.com` | toolspec |
| fomo | api_key | sealed | `api.fomo.com` | toolspec |
| formcarry | api_key | sealed | `formcarry.com` | toolspec (unverified) |
| formsite | api_key | sealed | `*.formsite.com` | toolspec |
| foursquare | api_key | sealed | `api.foursquare.com`, `places-api.foursquare.com` | toolspec |
| freshbooks | oauth2 | sealed | `api.freshbooks.com`, `auth.freshbooks.com` | toolspec |
| freshdesk | basic | sealed | `*.freshdesk.com` | toolspec |
| gmail | oauth2 | sealed | `gmail.googleapis.com` | toolspec |
| gong | oauth2 | sealed | `api.gong.io`, `*.api.gong.io` | toolspec |
| google-analytics | oauth2 | sealed | `analyticsdata.googleapis.com`, `analyticsadmin.googleapis.com`, `oauth2.googleapis.com` | toolspec |
| google-maps | api_key | sealed | `geocode.googleapis.com`, `places.googleapis.com`, `routes.googleapis.com` | toolspec |
| googleads | oauth2 | sealed | `googleads.googleapis.com`, `oauth2.googleapis.com` | toolspec |
| googlebigquery | oauth2 | sealed | `bigquery.googleapis.com` | toolspec |
| googlecalendar | oauth2 | sealed | `www.googleapis.com` | toolspec |
| googledocs | oauth2 | sealed | `docs.googleapis.com`, `www.googleapis.com`, `oauth2.googleapis.com` | toolspec |
| googledrive | oauth2 | sealed | `www.googleapis.com`, `oauth2.googleapis.com` | toolspec |
| googlemeet | oauth2 | sealed | `meet.googleapis.com`, `oauth2.googleapis.com` | toolspec |
| googlephotos | oauth2 | sealed | `photoslibrary.googleapis.com` | toolspec |
| googlesheets | oauth2 | sealed | `sheets.googleapis.com` | toolspec |
| googletasks | oauth2 | sealed | `tasks.googleapis.com`, `www.googleapis.com` | toolspec |
| gorgias | oauth2 | sealed | `gorgias.com`, `*.gorgias.com` | toolspec |
| gumroad | oauth2 | sealed | `api.gumroad.com` | toolspec |
| hackernews |  | sealed | `hacker-news.firebaseio.com`, `hn.algolia.com` | toolspec |
| hackerrank-work | api_key | sealed | `www.hackerrank.com` | toolspec |
| harvest | oauth2 | sealed | `api.harvestapp.com`, `id.getharvest.com` | toolspec |
| heygen | api_key | sealed | `api.heygen.com`, `upload.heygen.com` | toolspec |
| highlevel | oauth2 | sealed | `services.leadconnectorhq.com` | toolspec |
| hubspot | oauth2 | sealed | `api.hubapi.com` | toolspec |
| humanloop | api_key | sealed | `api.humanloop.com` | toolspec |
| intercom | oauth2 | sealed | `api.intercom.io` | toolspec |
| interzoid | api_key | sealed | `api.interzoid.com` | toolspec |
| jira | oauth2 | sealed | `api.atlassian.com`, `*.atlassian.net` | toolspec |
| junglescout | api_key | sealed | `developer.junglescout.com` | toolspec |
| klaviyo | oauth2 | sealed | `a.klaviyo.com` | toolspec |
| klipfolio | api_key | sealed | `app.klipfolio.com` | toolspec |
| kommo | oauth2 | sealed | `kommo.com`, `*.kommo.com` | toolspec |
| launchdarkly | api_key | sealed | `app.launchdarkly.com` | toolspec |
| lever | oauth2 | sealed | `api.lever.co`, `auth.lever.co` | toolspec |
| lexoffice | api_key | sealed | `api.lexware.io` | toolspec |
| linear | oauth2 | sealed | `api.linear.app` | toolspec |
| linkedin | oauth2 | sealed | `api.linkedin.com` | toolspec |
| linkhut | oauth2 | sealed | `api.ln.ht`, `ln.ht` | toolspec |
| linkup | api_key | sealed | `api.linkup.so` | toolspec |
| listennotes | api_key | sealed | `listen-api.listennotes.com` | toolspec |
| lmnt | api_key | sealed | `api.lmnt.com` | toolspec |
| mailchimp | oauth2 | sealed | `login.mailchimp.com`, `*.api.mailchimp.com` | toolspec |
| mailerlite | api_key | sealed | `connect.mailerlite.com` | toolspec (unverified) |
| maintainx | api_key | sealed | `api.getmaintainx.com` | toolspec |
| mem0 | api_key | sealed | `api.mem0.ai` | toolspec |
| metaads | oauth2 | sealed | `graph.facebook.com` | toolspec |
| metatextai | api_key | sealed | `api.metatext.ai`, `guard-api.metatext.ai` | toolspec |
| microsoft-clarity | api_key | sealed | `www.clarity.ms` | toolspec |
| microsoft-teams | oauth2 | sealed | `graph.microsoft.com`, `login.microsoftonline.com` | toolspec |
| miro | oauth2 | sealed | `api.miro.com` | toolspec |
| mixpanel | basic | sealed | `mixpanel.com`, `api.mixpanel.com`, `data.mixpanel.com`, `eu.mixpanel.com` | toolspec |
| mocean | api_key | entrusted | `rest.moceanapi.com` | toolspec |
| monday | oauth2 | sealed | `api.monday.com` | toolspec |
| mopinion | api_key | sealed | `api.mopinion.com` | toolspec |
| more-trees | api_key | sealed | `*.platform.moretrees.eco` | toolspec |
| moz | api_key | sealed | `api.moz.com`, `lsapi.seomoz.com` | toolspec |
| mural | oauth2 | sealed | `app.mural.co` | toolspec |
| neon | api_key | sealed | `console.neon.tech` | toolspec |
| netsuite | oauth2 | sealed | `*.suitetalk.api.netsuite.com`, `*.app.netsuite.com`, `system.netsuite.com` | toolspec |
| ngrok | api_key | sealed | `api.ngrok.com` | toolspec |
| notion | oauth2 | sealed | `api.notion.com` | toolspec |
| onedrive | oauth2 | sealed | `graph.microsoft.com` | toolspec |
| opensea | api_key | sealed | `api.opensea.io` | toolspec |
| outlook | oauth2 | sealed | `graph.microsoft.com`, `login.microsoftonline.com` | toolspec |
| pagerduty | oauth2 | sealed | `api.pagerduty.com` | toolspec |
| pandadoc | api_key | sealed | `api.pandadoc.com` | toolspec |
| peopledatalabs | api_key | sealed | `api.peopledatalabs.com` | toolspec |
| perplexityai | api_key | sealed | `api.perplexity.ai` | toolspec |
| piggy | api_key | sealed | `api.piggy.eu` | toolspec |
| pipedrive | oauth2 | sealed | `api.pipedrive.com`, `*.pipedrive.com` | toolspec |
| placekey | api_key | sealed | `api.placekey.io` | toolspec (unverified) |
| posthog | api_key | sealed | `app.posthog.com`, `us.posthog.com`, `eu.posthog.com`, `us.i.posthog.com`, `eu.i.posthog.com` | toolspec |
| process-street | api_key | sealed | `public-api.process.st` | toolspec |
| productboard | oauth2 | sealed | `api.productboard.com` | toolspec |
| rafflys | api_key | sealed | `app-sorteos.com` | toolspec (unverified) |
| recallai | api_key | sealed | `*.recall.ai` | toolspec |
| reddit | oauth2 | sealed | `oauth.reddit.com`, `www.reddit.com` | toolspec |
| retellai | api_key | sealed | `api.retellai.com` | toolspec |
| rocketlane | api_key | sealed | `api.rocketlane.com` | toolspec |
| rocketreach | api_key | sealed | `api.rocketreach.co` | toolspec |
| salesforce | oauth2 | sealed | `login.salesforce.com`, `*.salesforce.com`, `*.my.salesforce.com` | toolspec |
| screenshotone | api_key | entrusted | `api.screenshotone.com` | toolspec |
| semanticscholar | api_key | sealed | `api.semanticscholar.org` | toolspec |
| semrush | api_key | entrusted | `api.semrush.com` | toolspec |
| sendgrid | api_key | sealed | `api.sendgrid.com` | toolspec |
| sentry | oauth2 | sealed | `sentry.io`, `*.sentry.io` | toolspec |
| serpapi | api_key | entrusted | `serpapi.com` | toolspec |
| servicem8 | oauth2 | sealed | `api.servicem8.com` | toolspec |
| servicenow | oauth2 | sealed | `*.service-now.com` | toolspec |
| sharepoint | oauth2 | sealed | `graph.microsoft.com`, `*.sharepoint.com` | toolspec |
| shopify | oauth2 | sealed | `*.myshopify.com` | toolspec |
| shortcut | api_key | sealed | `api.app.shortcut.com` | toolspec |
| simplesat | api_key | sealed | `api.simplesat.io` | toolspec |
| slack | oauth2 | sealed | `slack.com`, `edgeapi.slack.com` | adopts github.com/korotovsky/slack-mcp-server v1.3.0 |
| slackbot | oauth2 | sealed | `slack.com`, `*.slack.com` | adopts github.com/korotovsky/slack-mcp-server v1.3.0 |
| smugmug | api_key | entrusted | `api.smugmug.com` | toolspec |
| snowflake | oauth2 | sealed | `*.snowflakecomputing.com`, `status.snowflake.com` | toolspec |
| square | oauth2 | sealed | `connect.squareup.com` | toolspec |
| stack-exchange | oauth2 | entrusted | `api.stackexchange.com` | toolspec |
| stripe | api_key | sealed | `api.stripe.com` | toolspec |
| supabase | api_key | sealed | `api.supabase.com` | toolspec |
| surveymonkey | oauth2 | sealed | `api.surveymonkey.com` | toolspec |
| tavily | api_key | sealed | `api.tavily.com` | toolspec |
| text-to-pdf | api_key | sealed | `v2.convertapi.com` | toolspec |
| textrazor | api_key | sealed | `api.textrazor.com` | toolspec |
| timecamp | api_key | sealed | `app.timecamp.com` | toolspec |
| timely | oauth2 | sealed | `api.timelyapp.com` | toolspec |
| tinypng | basic | sealed | `api.tinify.com` | toolspec |
| tinyurl | api_key | sealed | `api.tinyurl.com` | toolspec |
| tisane | api_key | sealed | `api.tisane.ai` | toolspec |
| todoist | oauth2 | sealed | `api.todoist.com`, `todoist.com` | toolspec |
| toneden | oauth2 | sealed | `www.toneden.io` | toolspec |
| trello | api_key | entrusted | `api.trello.com`, `trello.com` | toolspec |
| twitter | oauth2 | sealed | `api.twitter.com`, `api.x.com`, `upload.twitter.com` | toolspec |
| typefully | api_key | sealed | `api.typefully.com` | toolspec |
| waboxapp | api_key | entrusted | `www.waboxapp.com` | planned (no public API found) |
| wakatime | oauth2 | sealed | `api.wakatime.com`, `wakatime.com` | toolspec |
| weathermap | api_key | entrusted | `api.openweathermap.org` | toolspec |
| webex | oauth2 | sealed | `webexapis.com` | toolspec |
| webflow | oauth2 | sealed | `api.webflow.com` | toolspec |
| whatsapp | oauth2 | sealed | `graph.facebook.com` | toolspec |
| workiom | api_key | sealed | `workiom.com`, `*.workiom.com` | toolspec (unverified) |
| wrike | oauth2 | sealed | `www.wrike.com`, `*.wrike.com` | toolspec |
| xero | oauth2 | sealed | `api.xero.com`, `identity.xero.com` | toolspec |
| yandex | oauth2 | sealed | `cloud-api.yandex.net`, `api.music.yandex.net`, `api-metrika.yandex.net` | toolspec |
| ynab | oauth2 | sealed | `api.ynab.com`, `app.ynab.com` | toolspec (unverified) |
| yousearch | api_key | sealed | `ydc-index.io`, `api.ydc-index.io`, `chat-api.you.com`, `api.you.com` | toolspec |
| youtube | oauth2 | sealed | `www.googleapis.com`, `youtube.googleapis.com` | toolspec |
| zendesk | oauth2 | sealed | `*.zendesk.com` | toolspec |
| zenrows | api_key | entrusted | `api.zenrows.com` | toolspec |
| zenserp | api_key | sealed | `app.zenserp.com` | toolspec |
| zoho | oauth2 | sealed | `www.zohoapis.com`, `*.zohoapis.com`, `accounts.zoho.com` | toolspec |
| zoho-bigin | oauth2 | sealed | `www.zohoapis.com`, `*.zohoapis.com`, `accounts.zoho.com` | toolspec |
| zoho-books | oauth2 | sealed | `www.zohoapis.com`, `*.zohoapis.com`, `accounts.zoho.com` | toolspec |
| zoho-desk | oauth2 | sealed | `desk.zoho.com`, `accounts.zoho.com` | toolspec |
| zoho-inventory | oauth2 | sealed | `www.zohoapis.com`, `*.zohoapis.com`, `accounts.zoho.com` | toolspec |
| zoho-invoice | oauth2 | sealed | `www.zohoapis.com`, `accounts.zoho.com` | toolspec |
| zoho-mail | oauth2 | sealed | `mail.zoho.com`, `accounts.zoho.com` | toolspec |
| zoom | oauth2 | sealed | `api.zoom.us` | toolspec |
| zoominfo | oauth2 | sealed | `api.zoominfo.com` | toolspec |

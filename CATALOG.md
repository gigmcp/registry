# GigMCP Manifest Catalog

This catalog mirrors [Composio's toolkit list](https://composio.dev/tools) as a curated, lint-enforced set of MCP-server manifests. All entries except `echo` and `fetch` are **planned**: their manifests are hand-curated and pass `registryctl lint` (schema, tier, credential-injection, and egress-denylist checks), but the server implementations are pending and the image digests are placeholders — placeholder digests are **not installable**. The reference-implementation entries have their source in separate repos (see Status column).

| Name | Auth | Tier | Egress | Status |
|------|------|------|--------|--------|
| echo | api_key | sealed | `api.example.com` | reference implementation: github.com/gigmcp/gigmcp |
| fetch | api_key | sealed | `example.com`, `*.example.com` | reference implementation: github.com/gigmcp/gigmcp |
| ably | api_key | sealed | `rest.ably.io` | planned |
| acculynx | api_key | sealed | `api.acculynx.com` | planned |
| activecampaign | api_key | sealed | `*.api-us1.com` | planned |
| affinity | api_key | sealed | `api.affinity.co` | planned |
| agencyzoom | api_key | sealed | `api.agencyzoom.com` | planned |
| ahrefs | api_key | sealed | `api.ahrefs.com` | planned |
| airtable | oauth2 | sealed | `api.airtable.com` | planned |
| alchemy | api_key | entrusted | `*.g.alchemy.com` | planned |
| altoviz | api_key | sealed | `api.altoviz.com` | planned |
| amcards | api_key | sealed | `amcards.com` | planned |
| amplitude | basic | sealed | `amplitude.com`, `api2.amplitude.com` | planned |
| apaleo | oauth2 | sealed | `api.apaleo.com`, `identity.apaleo.com` | planned |
| apollo | api_key | sealed | `api.apollo.io` | planned |
| appdrag | api_key | sealed | `api.appdrag.com` | planned |
| asana | oauth2 | sealed | `app.asana.com` | planned |
| ashby | basic | sealed | `api.ashbyhq.com` | planned |
| attio | oauth2 | sealed | `api.attio.com` | planned |
| bamboohr | basic | sealed | `api.bamboohr.com` | planned |
| bannerbear | api_key | sealed | `api.bannerbear.com` | planned |
| baserow | api_key | sealed | `api.baserow.io` | planned |
| beeminder | api_key | entrusted | `www.beeminder.com` | planned |
| bitbucket | oauth2 | sealed | `api.bitbucket.org`, `bitbucket.org` | planned |
| bitwarden | oauth2 | sealed | `api.bitwarden.com`, `identity.bitwarden.com` | planned |
| blackbaud | oauth2 | sealed | `api.sky.blackbaud.com`, `oauth2.sky.blackbaud.com` | planned |
| blackboard | oauth2 | sealed | `*.blackboard.com` | planned |
| boldsign | oauth2 | sealed | `api.boldsign.com` | planned |
| bolna | api_key | sealed | `api.bolna.ai` | planned |
| borneo | api_key | sealed | `api.borneo.io`, `*.borneo.io` | planned |
| botbaba | api_key | sealed | `botbaba.io`, `*.botbaba.io` | planned |
| box | oauth2 | sealed | `api.box.com`, `upload.box.com` | planned |
| brandfetch | api_key | sealed | `api.brandfetch.io`, `graphql.brandfetch.io` | planned |
| breezy-hr | api_key | sealed | `api.breezy.hr` | planned |
| brevo | api_key | sealed | `api.brevo.com` | planned |
| brex | oauth2 | sealed | `platform.brexapis.com` | planned |
| browseai | api_key | sealed | `api.browse.ai` | planned |
| browserbase | api_key | sealed | `api.browserbase.com` | planned |
| bubble | api_key | sealed | `bubble.io`, `*.bubbleapps.io` | planned |
| cal | api_key | sealed | `api.cal.com` | planned |
| calendly | oauth2 | sealed | `api.calendly.com`, `auth.calendly.com` | planned |
| canva | oauth2 | sealed | `api.canva.com` | planned |
| canvas | api_key | sealed | `canvas.instructure.com`, `*.instructure.com` | planned |
| chatwork | api_key | sealed | `api.chatwork.com` | planned |
| chmeetings | api_key | sealed | `api.chmeetings.com`, `chmeetings.com` | planned |
| clickup | oauth2 | sealed | `api.clickup.com` | planned |
| close | basic | sealed | `api.close.com` | planned |
| cloudflare | api_key | sealed | `api.cloudflare.com` | planned |
| coda | api_key | sealed | `coda.io` | planned |
| coinbase | api_key | sealed | `api.coinbase.com`, `api.exchange.coinbase.com` | planned |
| coinmarketcal | api_key | sealed | `developers.coinmarketcal.com` | planned |
| confluence | oauth2 | sealed | `api.atlassian.com`, `*.atlassian.net` | planned |
| contentful | oauth2 | sealed | `api.contentful.com`, `cdn.contentful.com` | planned |
| crustdata | api_key | sealed | `api.crustdata.com` | planned |
| d2lbrightspace | oauth2 | sealed | `auth.brightspace.com`, `*.brightspace.com` | planned |
| dailybot | api_key | sealed | `api.dailybot.com` | planned |
| datadog | api_key | sealed | `api.datadoghq.com`, `*.datadoghq.com` | planned |
| datagma | api_key | entrusted | `gateway.datagma.net` | planned |
| datarobot | api_key | sealed | `app.datarobot.com` | planned |
| demio | api_key | sealed | `my.demio.com` | planned |
| dialpad | oauth2 | sealed | `dialpad.com` | planned |
| digicert | api_key | sealed | `www.digicert.com`, `api.digicert.com` | planned |
| discord | oauth2 | sealed | `discord.com` | planned |
| discordbot | api_key | sealed | `discord.com`, `cdn.discordapp.com` | planned |
| docmosis | api_key | sealed | `*.dws4.docmosis.com`, `*.dws3.docmosis.com` | planned |
| docusign | oauth2 | sealed | `account.docusign.com`, `*.docusign.net` | planned |
| dropbox | oauth2 | sealed | `api.dropboxapi.com`, `content.dropboxapi.com` | planned |
| dropbox-sign | oauth2 | sealed | `api.hellosign.com` | planned |
| dynamics365 | oauth2 | sealed | `*.dynamics.com`, `login.microsoftonline.com` | planned |
| echtpost | api_key | sealed | `api.echtpost.de` | planned |
| elevenlabs | api_key | sealed | `api.elevenlabs.io` | planned |
| eventbrite | oauth2 | sealed | `www.eventbriteapi.com` | planned |
| exa | api_key | sealed | `api.exa.ai` | planned |
| exist | oauth2 | sealed | `exist.io` | planned |
| facebook | oauth2 | sealed | `graph.facebook.com` | planned |
| figma | oauth2 | sealed | `api.figma.com` | planned |
| finage | api_key | entrusted | `api.finage.co.uk` | planned |
| firecrawl | api_key | sealed | `api.firecrawl.dev` | planned |
| fireflies | api_key | sealed | `api.fireflies.ai` | planned |
| flutterwave | api_key | sealed | `api.flutterwave.com` | planned |
| fomo | api_key | sealed | `api.fomo.com` | planned |
| formcarry | api_key | sealed | `formcarry.com` | planned |
| formsite | api_key | sealed | `*.formsite.com` | planned |
| foursquare | api_key | sealed | `api.foursquare.com`, `places-api.foursquare.com` | planned |
| freshbooks | oauth2 | sealed | `api.freshbooks.com`, `auth.freshbooks.com` | planned |
| freshdesk | basic | sealed | `*.freshdesk.com` | planned |
| gmail | oauth2 | sealed | `gmail.googleapis.com` | planned |
| gong | oauth2 | sealed | `api.gong.io`, `*.api.gong.io` | planned |
| google-analytics | oauth2 | sealed | `analyticsdata.googleapis.com`, `analyticsadmin.googleapis.com`, `oauth2.googleapis.com` | planned |
| google-maps | api_key | sealed | `maps.googleapis.com`, `places.googleapis.com`, `routes.googleapis.com` | planned |
| googleads | oauth2 | sealed | `googleads.googleapis.com`, `oauth2.googleapis.com` | planned |
| googlebigquery | oauth2 | sealed | `bigquery.googleapis.com` | planned |
| googlecalendar | oauth2 | sealed | `www.googleapis.com` | planned |
| googledocs | oauth2 | sealed | `docs.googleapis.com`, `www.googleapis.com`, `oauth2.googleapis.com` | planned |
| googledrive | oauth2 | sealed | `www.googleapis.com`, `oauth2.googleapis.com` | planned |
| googlemeet | oauth2 | sealed | `meet.googleapis.com`, `oauth2.googleapis.com` | planned |
| googlephotos | oauth2 | sealed | `photoslibrary.googleapis.com` | planned |
| googlesheets | oauth2 | sealed | `sheets.googleapis.com` | planned |
| googletasks | oauth2 | sealed | `tasks.googleapis.com`, `www.googleapis.com` | planned |
| gorgias | oauth2 | sealed | `gorgias.com`, `*.gorgias.com` | planned |
| gumroad | oauth2 | sealed | `api.gumroad.com` | planned |
| hackernews | none | sealed | `hacker-news.firebaseio.com`, `hn.algolia.com` | planned |
| hackerrank-work | api_key | sealed | `www.hackerrank.com` | planned |
| harvest | oauth2 | sealed | `api.harvestapp.com`, `id.getharvest.com` | planned |
| heygen | api_key | sealed | `api.heygen.com`, `upload.heygen.com` | planned |
| highlevel | oauth2 | sealed | `services.leadconnectorhq.com` | planned |
| hubspot | oauth2 | sealed | `api.hubapi.com` | planned |
| humanloop | api_key | sealed | `api.humanloop.com` | planned |
| intercom | oauth2 | sealed | `api.intercom.io` | planned |
| interzoid | api_key | sealed | `api.interzoid.com` | planned |
| jira | oauth2 | sealed | `api.atlassian.com`, `*.atlassian.net` | planned |
| junglescout | api_key | sealed | `developer.junglescout.com` | planned |
| klaviyo | oauth2 | sealed | `a.klaviyo.com` | planned |
| klipfolio | api_key | sealed | `app.klipfolio.com` | planned |
| kommo | oauth2 | sealed | `kommo.com`, `*.kommo.com` | planned |
| launchdarkly | api_key | sealed | `app.launchdarkly.com` | planned |
| lever | oauth2 | sealed | `api.lever.co`, `auth.lever.co` | planned |
| lexoffice | api_key | sealed | `api.lexoffice.io` | planned |
| linear | oauth2 | sealed | `api.linear.app` | planned |
| linkedin | oauth2 | sealed | `api.linkedin.com` | planned |
| linkhut | oauth2 | sealed | `api.ln.ht`, `ln.ht` | planned |
| linkup | api_key | sealed | `api.linkup.so` | planned |
| listennotes | api_key | sealed | `listen-api.listennotes.com` | planned |
| lmnt | api_key | sealed | `api.lmnt.com` | planned |
| mailchimp | oauth2 | sealed | `login.mailchimp.com`, `*.api.mailchimp.com` | planned |
| mailerlite | api_key | sealed | `connect.mailerlite.com` | planned |
| maintainx | api_key | sealed | `api.getmaintainx.com` | planned |
| mem0 | api_key | sealed | `api.mem0.ai` | planned |
| metaads | oauth2 | sealed | `graph.facebook.com` | planned |
| metatextai | api_key | sealed | `api.metatext.ai` | planned |
| microsoft-clarity | api_key | sealed | `www.clarity.ms` | planned |
| microsoft-teams | oauth2 | sealed | `graph.microsoft.com`, `login.microsoftonline.com` | planned |
| miro | oauth2 | sealed | `api.miro.com` | planned |
| mixpanel | basic | sealed | `mixpanel.com`, `api.mixpanel.com`, `data.mixpanel.com`, `eu.mixpanel.com` | planned |
| mocean | api_key | entrusted | `rest.moceanapi.com` | planned |
| monday | oauth2 | sealed | `api.monday.com` | planned |
| mopinion | api_key | sealed | `api.mopinion.com` | planned |
| more-trees | api_key | sealed | `api.moretrees.eco` | planned |
| moz | api_key | sealed | `api.moz.com`, `lsapi.seomoz.com` | planned |
| mural | oauth2 | sealed | `app.mural.co` | planned |
| neon | api_key | sealed | `console.neon.tech` | planned |
| netsuite | oauth2 | sealed | `*.suitetalk.api.netsuite.com`, `*.app.netsuite.com`, `system.netsuite.com` | planned |
| ngrok | api_key | sealed | `api.ngrok.com` | planned |
| notion | oauth2 | sealed | `api.notion.com` | planned |
| onedrive | oauth2 | sealed | `graph.microsoft.com` | planned |
| opensea | api_key | sealed | `api.opensea.io` | planned |
| outlook | oauth2 | sealed | `graph.microsoft.com`, `login.microsoftonline.com` | planned |
| pagerduty | oauth2 | sealed | `api.pagerduty.com` | planned |
| pandadoc | api_key | sealed | `api.pandadoc.com` | planned |
| peopledatalabs | api_key | sealed | `api.peopledatalabs.com` | planned |
| perplexityai | api_key | sealed | `api.perplexity.ai` | planned |
| piggy | api_key | sealed | `api.piggy.eu` | planned |
| pipedrive | oauth2 | sealed | `api.pipedrive.com`, `*.pipedrive.com` | planned |
| placekey | api_key | sealed | `api.placekey.io` | planned |
| posthog | api_key | sealed | `app.posthog.com`, `us.posthog.com`, `eu.posthog.com` | planned |
| process-street | api_key | sealed | `public-api.process.st` | planned |
| productboard | oauth2 | sealed | `api.productboard.com` | planned |
| rafflys | api_key | sealed | `api.app-sorteos.com` | planned |
| recallai | api_key | sealed | `*.recall.ai` | planned |
| reddit | oauth2 | sealed | `oauth.reddit.com`, `www.reddit.com` | planned |
| retellai | api_key | sealed | `api.retellai.com` | planned |
| rocketlane | api_key | sealed | `api.rocketlane.com` | planned |
| rocketreach | api_key | sealed | `api.rocketreach.co` | planned |
| salesforce | oauth2 | sealed | `login.salesforce.com`, `*.salesforce.com`, `*.my.salesforce.com` | planned |
| screenshotone | api_key | entrusted | `api.screenshotone.com` | planned |
| semanticscholar | api_key | sealed | `api.semanticscholar.org` | planned |
| semrush | api_key | entrusted | `api.semrush.com` | planned |
| sendgrid | api_key | sealed | `api.sendgrid.com` | planned |
| sentry | oauth2 | sealed | `sentry.io`, `*.sentry.io` | planned |
| serpapi | api_key | entrusted | `serpapi.com` | planned |
| servicem8 | oauth2 | sealed | `api.servicem8.com` | planned |
| servicenow | oauth2 | sealed | `*.service-now.com` | planned |
| sharepoint | oauth2 | sealed | `graph.microsoft.com`, `*.sharepoint.com` | planned |
| shopify | oauth2 | sealed | `*.myshopify.com` | planned |
| shortcut | api_key | sealed | `api.app.shortcut.com` | planned |
| simplesat | api_key | sealed | `api.simplesat.io` | planned |
| slack | oauth2 | sealed | `slack.com`, `*.slack.com` | planned |
| slackbot | oauth2 | sealed | `slack.com`, `*.slack.com` | planned |
| smugmug | api_key | entrusted | `api.smugmug.com` | planned |
| snowflake | oauth2 | sealed | `*.snowflakecomputing.com`, `status.snowflake.com` | planned |
| square | oauth2 | sealed | `connect.squareup.com` | planned |
| stack-exchange | oauth2 | entrusted | `api.stackexchange.com` | planned |
| stripe | api_key | sealed | `api.stripe.com` | planned |
| supabase | api_key | sealed | `api.supabase.com` | planned |
| surveymonkey | oauth2 | sealed | `api.surveymonkey.com` | planned |
| tavily | api_key | sealed | `api.tavily.com` | planned |
| text-to-pdf | api_key | sealed | `v2.convertapi.com` | planned |
| textrazor | api_key | sealed | `api.textrazor.com` | planned |
| timecamp | api_key | sealed | `app.timecamp.com` | planned |
| timely | oauth2 | sealed | `api.timelyapp.com` | planned |
| tinypng | basic | sealed | `api.tinify.com` | planned |
| tinyurl | api_key | sealed | `api.tinyurl.com` | planned |
| tisane | api_key | sealed | `api.tisane.ai` | planned |
| todoist | oauth2 | sealed | `api.todoist.com`, `todoist.com` | planned |
| toneden | oauth2 | sealed | `api.toneden.io` | planned |
| trello | api_key | entrusted | `api.trello.com`, `trello.com` | planned |
| twitter | oauth2 | sealed | `api.twitter.com`, `api.x.com`, `upload.twitter.com` | planned |
| typefully | api_key | sealed | `api.typefully.com` | planned |
| waboxapp | api_key | entrusted | `www.waboxapp.com` | planned |
| wakatime | oauth2 | sealed | `api.wakatime.com`, `wakatime.com` | planned |
| weathermap | api_key | entrusted | `api.openweathermap.org` | planned |
| webex | oauth2 | sealed | `webexapis.com` | planned |
| webflow | oauth2 | sealed | `api.webflow.com` | planned |
| whatsapp | oauth2 | sealed | `graph.facebook.com` | planned |
| workiom | api_key | sealed | `workiom.com`, `*.workiom.com` | planned |
| wrike | oauth2 | sealed | `www.wrike.com`, `*.wrike.com` | planned |
| xero | oauth2 | sealed | `api.xero.com`, `identity.xero.com` | planned |
| yandex | oauth2 | sealed | `cloud-api.yandex.net`, `api.music.yandex.net`, `api-metrika.yandex.net` | planned |
| ynab | oauth2 | sealed | `api.ynab.com`, `app.ynab.com` | planned |
| yousearch | api_key | sealed | `api.ydc-index.io`, `chat-api.you.com` | planned |
| youtube | oauth2 | sealed | `www.googleapis.com`, `youtube.googleapis.com` | planned |
| zendesk | oauth2 | sealed | `*.zendesk.com` | planned |
| zenrows | api_key | entrusted | `api.zenrows.com` | planned |
| zenserp | api_key | sealed | `app.zenserp.com` | planned |
| zoho | oauth2 | sealed | `www.zohoapis.com`, `*.zohoapis.com`, `accounts.zoho.com` | planned |
| zoho-bigin | oauth2 | sealed | `www.zohoapis.com`, `*.zohoapis.com`, `accounts.zoho.com` | planned |
| zoho-books | oauth2 | sealed | `www.zohoapis.com`, `*.zohoapis.com`, `accounts.zoho.com` | planned |
| zoho-desk | oauth2 | sealed | `desk.zoho.com`, `accounts.zoho.com` | planned |
| zoho-inventory | oauth2 | sealed | `www.zohoapis.com`, `*.zohoapis.com`, `accounts.zoho.com` | planned |
| zoho-invoice | oauth2 | sealed | `www.zohoapis.com`, `accounts.zoho.com` | planned |
| zoho-mail | oauth2 | sealed | `mail.zoho.com`, `accounts.zoho.com` | planned |
| zoom | oauth2 | sealed | `api.zoom.us` | planned |
| zoominfo | oauth2 | sealed | `api.zoominfo.com` | planned |

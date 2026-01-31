# XMLTV Wiki

> Source: https://wiki.xmltv.org/index.php/XMLTVFormat
> Fetched: 2026-01-31T15:59:42.882468+00:00
> Content-Hash: 6fc219c8e4593fe0
> Type: html

---

# XMLTVFormat

From XMLTV

Jump to: navigation, search

# XMLTV File format

The format used differs from most other XML-based TV listings formats in that it is written from the user's point of view, rather that the broadcaster's. It doesn't divide listings into channels, instead all the channels are mixed together into a single unified listing. Each programme has details such as name, description, and credits stored as supplements, but metadata like broadcast details are stored as attributes. There is support for listings in multiple languages and each programme can have 'language' and 'original language' details. 

The XMLTV File Format was originally created by Ed Avis, and is currently maintained by the [XMLTVProject](/index.php/XMLTVProject "XMLTVProject"). The current Git version of the DTD is available [here](https://github.com/XMLTV/xmltv/blob/master/xmltv.dtd). 

There are additional requirements on grabbers if they want to be "baseline compliant". See [XmltvCapabilities](/index.php/XmltvCapabilities "XmltvCapabilities")

Since the [DTD](https://github.com/XMLTV/xmltv/blob/master/xmltv.dtd) is available, you can also use [XmltvValidation](/index.php/XmltvValidation "XmltvValidation"). 

# Details

An XMLTV file has 2 types of records 

  * 'channel' records, store information about channels
  * 'program' records, store information about individual episodes



Most of the information is optional and may not be available from all sources 

This is what a sample xmltv file looks like 
    
    
    <?xml version="1.0" encoding="ISO-8859-1"?>
    <!DOCTYPE tv SYSTEM "xmltv.dtd">
    
    <tv source-info-url="http://www.schedulesdirect.org/" source-info-name="Schedules Direct" generator-info-name="XMLTV/$Id: tv_grab_na_dd.in,v 1.70 2008/03/03 15:21:41 rmeden Exp $" generator-info-url="http://www.xmltv.org/">
      <channel id="I10436.labs.zap2it.com">
        <display-name>13 KERA</display-name>
        <display-name>13 KERA TX42822:-</display-name>
        <display-name>13</display-name>
        <display-name>13 KERA fcc</display-name>
        <display-name>KERA</display-name>
        <display-name>KERA</display-name>
        <display-name>PBS Affiliate</display-name>
        <icon src="file://C:\Perl\site/share/xmltv/icons/KERA.gif" />
      </channel>
      <channel id="I10759.labs.zap2it.com">
        <display-name>11 KTVT</display-name>
        <display-name>11 KTVT TX42822:-</display-name>
        <display-name>11</display-name>
        <display-name>11 KTVT fcc</display-name>
        <display-name>KTVT</display-name>
        <display-name>KTVT</display-name>
        <display-name>CBS Affiliate</display-name>
        <icon src="file://C:\Perl\site/share/xmltv/icons/KTVT.gif" />
      </channel>
      <programme start="20080715003000 -0600" stop="20080715010000 -0600" channel="I10436.labs.zap2it.com">
        <title lang="en">NOW on PBS</title>
        <desc lang="en">Jordan's Queen Rania has made job creation a priority to help curb the staggering unemployment rates among youths in the Middle East.</desc>
        <date>20080711</date>
        <category lang="en">Newsmagazine</category>
        <category lang="en">Interview</category>
        <category lang="en">Public affairs</category>
        <category lang="en">Series</category>
        <episode-num system="dd_progid">EP01006886.0028</episode-num>
        <episode-num system="onscreen">427</episode-num>
        <audio>
          <stereo>stereo</stereo>
        </audio>
        <previously-shown start="20080711000000" />
        <subtitles type="teletext" />
      </programme>
      <programme start="20080715010000 -0600" stop="20080715023000 -0600" channel="I10436.labs.zap2it.com">
        <title lang="en">Mystery!</title>
        <sub-title lang="en">Foyle's War, Series IV: Bleak Midwinter</sub-title>
        <desc lang="en">Foyle investigates an explosion at a munitions factory, which he comes to believe may have been premeditated.</desc>
        <date>20070701</date>
        <category lang="en">Anthology</category>
        <category lang="en">Mystery</category>
        <category lang="en">Series</category>
        <episode-num system="dd_progid">EP00003026.0665</episode-num>
        <episode-num system="onscreen">2705</episode-num>
        <audio>
          <stereo>stereo</stereo>
        </audio>
        <previously-shown start="20070701000000" />
        <subtitles type="teletext" />
      </programme>
      <programme start="20080715023000 -0600" stop="20080715040000 -0600" channel="I10436.labs.zap2it.com">
        <title lang="en">Mystery!</title>
        <sub-title lang="en">Foyle's War, Series IV: Casualties of War</sub-title>
        <desc lang="en">The murder of a prominent scientist may have been due to a gambling debt.</desc>
        <date>20070708</date>
        <category lang="en">Anthology</category>
        <category lang="en">Mystery</category>
        <category lang="en">Series</category>
        <episode-num system="dd_progid">EP00003026.0666</episode-num>
        <episode-num system="onscreen">2706</episode-num>
        <audio>
          <stereo>stereo</stereo>
        </audio>
        <previously-shown start="20070708000000" />
        <subtitles type="teletext" />
      </programme>
      <programme start="20080715040000 -0600" stop="20080715043000 -0600" channel="I10436.labs.zap2it.com">
        <title lang="en">BBC World News</title>
        <desc lang="en">International issues.</desc>
        <category lang="en">News</category>
        <category lang="en">Series</category>
        <episode-num system="dd_progid">SH00315789.0000</episode-num>
        <previously-shown />
        <subtitles type="teletext" />
      </programme>
      <programme start="20080715043000 -0600" stop="20080715050000 -0600" channel="I10436.labs.zap2it.com">
        <title lang="en">Sit and Be Fit</title>
        <date>20070924</date>
        <category lang="en">Exercise</category>
        <category lang="en">Series</category>
        <episode-num system="dd_progid">EP00003847.0074</episode-num>
        <episode-num system="onscreen">901</episode-num>
        <audio>
          <stereo>stereo</stereo>
        </audio>
        <previously-shown start="20070924000000" />
        <subtitles type="teletext" />
      </programme>
      <programme start="20080715060000 -0600" stop="20080715080000 -0600" channel="I10759.labs.zap2it.com">
        <title lang="en">The Early Show</title>
        <desc lang="en">Republican candidate John McCain; premiere of the film "The Dark Knight."</desc>
        <date>20080715</date>
        <category lang="en">Talk</category>
        <category lang="en">News</category>
        <category lang="en">Series</category>
        <episode-num system="dd_progid">EP00337003.2361</episode-num>
        <audio>
          <stereo>stereo</stereo>
        </audio>
        <subtitles type="teletext" />
      </programme>
      <programme start="20080715080000 -0600" stop="20080715090000 -0600" channel="I10759.labs.zap2it.com">
        <title lang="en">Rachael Ray</title>
        <desc lang="en">Actresses Kim Raver, Brooke Shields and Lindsay Price ("Lipstick Jungle"); women in their 40s tell why they got breast implants; a 30-minute meal.</desc>
        <credits>
          <presenter>Rachael Ray</presenter>
        </credits>
        <date>20080306</date>
        <category lang="en">Talk</category>
        <category lang="en">Series</category>
        <episode-num system="dd_progid">EP00847333.0303</episode-num>
        <episode-num system="onscreen">2119</episode-num>
        <audio>
          <stereo>stereo</stereo>
        </audio>
        <previously-shown start="20080306000000" />
        <subtitles type="teletext" />
      </programme>
      <programme start="20080715090000 -0600" stop="20080715100000 -0600" channel="I10759.labs.zap2it.com">
        <title lang="en">The Price Is Right</title>
        <desc lang="en">Contestants bid for prizes then compete for fabulous showcases.</desc>
        <credits>
          <director>Bart Eskander</director>
          <producer>Roger Dobkowitz</producer>
          <presenter>Drew Carey</presenter>
        </credits>
        <category lang="en">Game show</category>
        <category lang="en">Series</category>
        <episode-num system="dd_progid">SH00004372.0000</episode-num>
        <audio>
          <stereo>stereo</stereo>
        </audio>
        <subtitles type="teletext" />
        <rating system="VCHIP">
          <value>TV-G</value>
        </rating>
      </programme>
      <programme start="20080715100000 -0600" stop="20080715103000 -0600" channel="I10759.labs.zap2it.com">
        <title lang="en">Jeopardy!</title>
        <credits>
          <presenter>Alex Trebek</presenter>
        </credits>
        <date>20080715</date>
        <category lang="en">Game show</category>
        <category lang="en">Series</category>
        <episode-num system="dd_progid">EP00002348.1700</episode-num>
        <episode-num system="onscreen">5507</episode-num>
        <subtitles type="teletext" />
        <rating system="VCHIP">
          <value>TV-G</value>
        </rating>
      </programme>
      <programme start="20080715103000 -0600" stop="20080715113000 -0600" channel="I10759.labs.zap2it.com">
        <title lang="en">The Young and the Restless</title>
        <sub-title lang="en">Sabrina Offers Victoria a Truce</sub-title>
        <desc lang="en">Jeff thinks Kyon stole the face cream; Nikki asks Nick to give David a chance; Amber begs Adrian to go to Australia.</desc>
        <credits>
          <actor>Peter Bergman</actor>
          <actor>Eric Braeden</actor>
          <actor>Jeanne Cooper</actor>
          <actor>Melody Thomas Scott</actor>
        </credits>
        <date>20080715</date>
        <category lang="en">Soap</category>
        <category lang="en">Series</category>
        <episode-num system="dd_progid">EP00004422.1359</episode-num>
        <episode-num system="onscreen">8937</episode-num>
        <audio>
          <stereo>stereo</stereo>
        </audio>
        <subtitles type="teletext" />
        <rating system="VCHIP">
          <value>TV-14</value>
        </rating>
      </programme>
    </tv>
    

Retrieved from "[https://wiki.xmltv.org/index.php?title=XMLTVFormat&oldid=2184](https://wiki.xmltv.org/index.php?title=XMLTVFormat&oldid=2184)" 

## Navigation menu

### Personal tools

  * [Log in](/index.php?title=Special:UserLogin&returnto=XMLTVFormat "You are encouraged to log in; however, it is not mandatory \[o\]")



### Namespaces

  * [Page](/index.php/XMLTVFormat "View the content page \[c\]")
  * [Discussion](/index.php?title=Talk:XMLTVFormat&action=edit&redlink=1 "Discussion about the content page \[t\]")



###  Variants




### Views

  * [Read](/index.php/XMLTVFormat)
  * [View source](/index.php?title=XMLTVFormat&action=edit "This page is protected.
You can view its source \[e\]")
  * [View history](/index.php?title=XMLTVFormat&action=history "Past revisions of this page \[h\]")



### More




###  Search

[](/index.php/Main_Page "Visit the main page")

### Helpful Links

  * [Main Page](/index.php/Main_Page)
  * [GitHub Project](https://github.com/XMLTV/)
  * [(previous) SF Project](http://sourceforge.net/projects/xmltv/)
  * [Download Info](/index.php/XMLTVProjectDownload)
  * [Wiki Changes](/index.php/Special:RecentChanges)



### Tools

  * [What links here](/index.php/Special:WhatLinksHere/XMLTVFormat "A list of all wiki pages that link here \[j\]")
  * [Related changes](/index.php/Special:RecentChangesLinked/XMLTVFormat "Recent changes in pages linked from this page \[k\]")
  * [Special pages](/index.php/Special:SpecialPages "A list of all special pages \[q\]")
  * [Printable version](/index.php?title=XMLTVFormat&printable=yes "Printable version of this page \[p\]")
  * [Permanent link](/index.php?title=XMLTVFormat&oldid=2184 "Permanent link to this revision of the page")
  * [Page information](/index.php?title=XMLTVFormat&action=info "More information about this page")



  * This page was last modified on 28 January 2019, at 03:54.


  * [Privacy policy](/index.php/XMLTV:Privacy_policy "XMLTV:Privacy policy")
  * [About XMLTV](/index.php/XMLTV:About "XMLTV:About")
  * [Disclaimers](/index.php/XMLTV:General_disclaimer "XMLTV:General disclaimer")


  * [](//www.mediawiki.org/)


  *[↑]: Back to Top

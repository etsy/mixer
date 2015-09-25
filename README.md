# Mixer

## Intro

Mixer is a tool to initiate meetings by randomly pairing individuals. Read more about it [here](https://codeascraft.com/2015/09/15/assisted-serendipity).

## Development Setup

0. [Setup](https://golang.org/doc/install) your golang environment

1. Grab the repo and put it in the right place:

        go get github.com/etsy/mixer
        cd $GOPATH/src/github.com/etsy/mixer

1. Create the 'mixer' database:

        mysql -u root
        create database mixer;
        \u mixer
        source db/mixer.sql
        source db/mixer_data.sql
        grant all privileges on mixer.* to 'mixeruser'@'127.0.0.1' identified by '<somepassword>';

1. Configure the mixer app:

  * `cp config.cfg.sample config.cfg` and edit the values to match your environment
  * `cp config-secrets.cfg.sample config-secrets.cfg` and put the password you created for `mixeruser` above in `config-secrets.cfg`

1. Assuming you're running apache, create a virtual host config for mixer in `/etc/httpd/conf.d/mixer.conf`:

        <Virtualhost *:80>
          ServerName mixer.<hostname>
          ProxyPass / http://localhost:3001/ retry=0
          ProxyPassReverse / http://localhost:3001/

          RequestHeader set X-Username <username>
        </VirtualHost>

**Note:** This is a development setup. In production, your authentication system needs to pass the username as an http header after the user has logged in. The header name can be defined in the config. This app was not written to assume any implementation of your authentication system.

1. Fire it up:

        go run mixer.go

1. Open it in your browser at http://mixer._<hostname_>

## Architecture

The server is written in Go, while the client-side javascript is backbone + jquery. There is a local *staff* table, which is intended to be an updated copy of data from your company-wide staff database. It is assumed in the example import script that your company provides an api to the staff database in JSON format which is defined below. If this is not the case you are welcome to write your own import implementation to populate the staff table.

Additionally, as mentioned in the setup, this assumes that the application is fronted by an authentication system that passes headers into the app that define the username of the logged-in user. This allows the application to see if the user has an account, or if they have to be presented with a registration form.

It is easy to add new mixers by inserting into the *groups* table. There is currently no user interface for this admin functionality. 

### Running the Mixer

Each mixer must be run via separate cronjobs. This gives some flexibility to the frequency each mixer is run. An example of what this cron entry would like like is below:

    # run as 12:00 on Mondays
    0 12 * * 1 cd /usr/local/mixer && /usr/local/mixer/bin/match --group=Engineering >> /var/log/mixer/engineeringcron.log 2>&1

### Importing the Staff Data

The local staff table is assumed to be populated and updated. There are some helper scripts written that make some assumptions. They assume that you are providing a URL returning JSON data of a specific format defined below. This is taken and inserted into the local staff table.

An example cron of running the import of staff data:

     # run script to sync staff data
     0 */4 * * * cd /usr/local/mixer && /usr/local/mixer/bin/staff >> /var/log/mixer/staffcron.log 2>&1

#### JSON Format of Staff Datafeed

    {
      id: 5000,
      auth_username: "flastname",
      first_name: "Test",
      last_name: "User",
      title: "Engineer",
      avatar: "https://someurl/img.jpg",
      enabled: 1,
      is_manager: false,
    }

### Cleaning up Alumni

Since it doesn't make sense to pair people with others that have left the company, some cleanup work is involved with the groups people have joined. There is a script to help here. The assumption is that you populate your staff table with both enabled and non-enabled users. The cleanup script will then remove any disabled users from the people_groups table so that they are no longer matched but still retain their account and history.

An example cron to run the cleanup script is:

        # run script to cleanup alumni users
        0 12 * * 1 cd /usr/local/mixer && /usr/local/mixer/bin/cleanup >> /var/log/mixer/cleanupcron.log 2>&1

## Notes

### the is_manager field

There is one specific mixer group called *Managers* that is limited to managers by some fairly specific code that should be abstracted at some point. This value is set by the *is_manager* field from your staff database.

### running alternating mixers

It is possible with some bash-fu to short circuit it to run every other week. Here is an example:

    # run engineering mixer every other monday starting from Jan 5 2015
    0 14 * * 1 bash -c '((($(date +\%s) + 259200) / 86400 \% 14))' && cd /usr/local/mixer && /usr/local/mixer/bin/match --group=Engineering >> /var/log/mixer/engineeringcron.log 2>&1
    # run manager mixer every other monday starting from Jan 12 2015
    0 14 * * 1 bash -c '((($(date +\%s) + 864000) / 86400 \% 14))' && cd /usr/local/mixer && /usr/local/mixer/bin/match --group=Managers >> /var/log/mixer/managerscron.log 2>&1


## TODO

See [Issues](https://github.com/etsy/mixer/issues) for bugs and ideas. We'd love to see your pull requests!

## Help

Stop by the [#codeascraft](irc://irc.freenode.net/codeascraft) IRC channel on Freenode if you need any help.


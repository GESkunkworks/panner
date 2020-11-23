# panner
Pan for EBS snapshot gold and maybe you'll get rich.

Analyzes the EBS snapshots in an account and determines how much money you could save by deleting snapshots that are missing their volumes. 

* Supports date filtering of snapshots
* determines whether snapshots have current volumes and outputs all to file
* outputs snapshots aggregated by common volume ID into a separate file so minimum estimated cost savings can be determined based on volume size.

This executable is powered by [dustcollector](https://github.com/GESkunkworks/dustcollector).

# Overview
EBS snapshots are billed based on GB-month rate. So if you have a 500GB volume with a single snapshot and the rate is $0.05 per GB-month then that snapshot will cost $25/month to store. Snapshots after the initial snapshot are only the difference of what's changed since the last snapshot. So if you snapshot the volume again and you changed 50GB worth of data then you are billed for 550 GB-month of data so your bill would increase to $27.50 the next month. 

If an EBS volume is deleted then the snapshot will remain in case you want to restore the volume at some point in the future. However, most of the time people just simply forget to delete snapshots when they terminate infrastructure or they intend to keep the snapshot for only a few months but then forget about it. This can leave snapshots laying around for years and the costs can add up. 

This tool is designed to help find that snapshot volume so you can clean it up and save some money. 

## Installation
Download a release from the releases tab for your OS architecture and unzip to a folder then execute the binary from command line. 

Alternatively if you have a Golang dev environment you can build locally with `make build`.

## Usage
The following command will analyze the EBS snapshots in the `digital-public-cloudops` account and output the snapshot info and cost results to two files `digital-public-cloudops-snapshots.csv` and `digital-public-cloudops-bars.csv` respectively. It will only look at snapshots that were created before `2019-01-01`. It will paginate snapshot results up to 20 pages with a page size of 750 (the max is 1000). 

```
$ ACCOUNT=digital-public-cloudops bash -c './panner -max-pages 20 -pagesize 750 -datefilter 2019-01-01 -profile $ACCOUNT -outfile-snapshots $ACCOUNT-snapshots.csv -outfile-bars $ACCOUNT-bars.csv -outfile-summary $ACCOUNT-summary.txt'
```

sample output:
```
t=2020-08-13T01:34:10-0400 lvl=info msg="Starting panner"
t=2020-08-13T01:34:10-0400 lvl=info msg="starting session" profile=digital-public-cloudops
t=2020-08-13T01:34:11-0400 lvl=info msg="Filtered snapshots page by date" pre-filter=13 post-filter=10 pageNum=1
t=2020-08-13T01:34:11-0400 lvl=info msg="searching for batch of volumes" size=8
t=2020-08-13T01:34:11-0400 lvl=info msg="Waiting for describeVolume batches to finish"
t=2020-08-13T01:34:14-0400 lvl=info msg="Total snapshots post date filter" snapshots_in_scope=10
t=2020-08-13T01:34:14-0400 lvl=info msg="Total snapshots analyzed" total-analyzed=13
t=2020-08-13T01:34:14-0400 lvl=info msg="grabbing all latest launch template versions"
t=2020-08-13T01:34:15-0400 lvl=info msg="Writing snapshots to file"
t=2020-08-13T01:34:15-0400 lvl=info msg="wrote nuggets to file" filename=digital-public-cloudops-snapshots.csv
t=2020-08-13T01:34:15-0400 lvl=info msg="Writing cost info to file"
t=2020-08-13T01:34:15-0400 lvl=info msg="wrote bars to file" filename=digital-public-cloudops-bars.csv
t=2020-08-13T01:34:15-0400 lvl=info msg="wrote summary to file" filename=digital-public-cloudops-summary.txt
```

From there you can look at the summary file and delete the snapshots if you want to realize the savings.  

Summary
```
After analyzing the account we can see that there are 5 snapshots that can be deleted because they were created before 2019-01-01 and are not used in any AutoScaling group or AMI sharing capacity. However, before these snapshots can be deleted several other resources need to be deleted first. Below you can find the ordered deletion plan:


Some of the snapshots we need to delete are currently registered as AMIs or used in Launch Templates/Configs. However we've detected that those AMI's and Launch Templates/Configs are not used in any autoscaling group. This doesn't mean they're not being used by someone (e.g., referenced in a cloudformation template). You should be safe to delete them but you should always check to be sure

If you feel comfortable then here's the plan:

Delete the following LaunchTemplates first:
        test-lt
then delete the following LaunchConfigurations:
        test-snap-lc
then delete the following AMIs:
        ami-a7ce9bdd
        ami-6cee4b16
then finally delete the following Snapshots:
        snap-092ab265885243a2d
        snap-005ccdfd0fedb77b6
        snap-06e70bf98b9e43b2f
        snap-0a4795e305f1bc40d
        snap-07a4f8539c10e0dc7
3 snapshots were spared because their EBS volume still exists
1 snapshots were spared because they were associated with an autoscaling group, were shared directly to another account, or were registered as an AMI that was shared to another account.
Total size of eligible for deletion is 40 GB. At a per GB-month rate of $0.050000 there is a potential savings of $2.000000
```

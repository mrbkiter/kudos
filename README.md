# kudos

Kudos is the simple and easy to use employee recognition software that enhances employee engagement and team communication. The first version (0.0.1) is integration with Slack chat which includes /kudos and /kudos-report commands

## Syntax

There are 3 supported commands: 

``/kudos Thanks @mrbkiter, @abc for great support``

This command would kudos mrbkiter and abc. It would increase +1 for each. There is kudos syntax you need to follow:

``/kudos (thanks|great|good|space*) (any words but not <>) <list_of_users> something_else_here``

For example: 

Accepted kudos:

``/kudos thanks @mrbkiter for helping ... ``

``/kudos thank you @mrbkiter. You did great job``

``/kudos great job @mrbkiter. Great release and helped @abc`` (counted for @mrbkiter)

``/kudos @mrbkiter @abc. You did great release``

Unaccepted kudos:

``/kudos abcbasdasd @mrbkiter for something`` 

``/kudos-report <THIS_MONTH | LAST_MONTH | THIS_WEEK | LAST_WEEK> @mrbkiter @abc ``

This command would return report for mrbkiter and abc (in case you input empty, it would return all user report). There are 4 report time types: THIS_MONTH, LAST_MONTH, THIS_WEEK, LAST_WEEK (default is THIS_WEEK) The result is returned in desc order

``/kudos-report detail @mrbkiter THIS_MONTH`` 

This command would return report detail for mrbkiter in this month. 

## The Architecture

The project is written in Golang, using aws dynamodb for db, aws API Gateway and Lambda. You can see the flow below: 

``slack --> AWS API Gateway --> Lambda --> Dynamodb ``

### Dynamodb design

The table need preconfigured partition key id1 and sort key id2 (both are string). Currently there are 2 main types: command and report. Whenever users type a kudos command, their command would be stored in ddb. (breaking down to each of kudos users to one row). We besides also calculate week number of the year and yyyy-mm for report builder. 

Example: 

``{
  "channelId": "C023G6N6D5E",
  "id1": "T022PA5N7KP#U022SCGDY58",
  "id2": "2232569221158.2091345755669.efa17f545015716af8bfa28f0ca96208",
  "msgId": "2232569221158.2091345755669.efa17f545015716af8bfa28f0ca96208",
  "sourceUserId": "U024U032H8A",
  "teamId": "T022PA5N7KP",
  "teamIdMonth": "T022PA5N7KP#2021-07",
  "teamIdWeek": "T022PA5N7KP#2021#26",
  "text": "/kudos you did great job <@U024D6VQX7Z|vu.yen.nguyen.88> <@U022SCGDY58|mrbkiter> <@U024U032H8A|vu.nguyen>",
  "timestamp": 1625411561,
  "ttl": 1633187561,
  "type": "**command**",
  "userId": "U022SCGDY58",
  "username": "mrbkiter"
}``

For report, it is triggered by ddb event trigger (which you need to enable at Trigger tab of your table). The trigger would be connected to a lambda function (report folder in this project) to help pre-calculate MONTH and WEEK total kudos of users. 

For example: 

Week report: 

``{
  "count": 3,
  "id1": "T022PA5N7KP#report",
  "id2": "2021#26#U024D6VQX7Z",
  "teamId": "T022PA5N7KP",
  "userId": "U024D6VQX7Z",
  "username": "vu.yen.nguyen.88"
}``

Monthly report: 

``{
  "count": 3,
  "id1": "T022PA5N7KP#report",
  "id2": "2021-07#U024D6VQX7Z",
  "teamId": "T022PA5N7KP",
  "userId": "U024D6VQX7Z",
  "username": "vu.yen.nguyen.88"
}``

If you notice, the partition key is composed as <team_id>#report (for report type), and <team_id>#<user_id> for command type. The sort key of report_type is <yyyy-MM>#<user_id>, or <yyyy>#<week_no>#<user_id> for week report. The id2 of command type is message id. 
  
If you need to extend your business, there are more rooms for you (we store teamId, channelId, ... so you can build more secondary index for your query)
  
### Code structure
  
  The code structure is simple: slack folder (for slack integration), report folder (for report lambda func). There are repos, model, ddb_entity which are for internal purposes (in case you need another db, just overwrite repo interface) 
  
### API Gateway configuration
  
  The GW would be configured to Slack lambda function. 
  
  ![image](https://user-images.githubusercontent.com/10323118/124391376-8d090980-dd1a-11eb-92be-c6510abe6ec9.png)

  ![image](https://user-images.githubusercontent.com/10323118/124391410-b5910380-dd1a-11eb-85ff-4cbc90d46b5e.png)

 ### Lambda Function
  
As explained above, you would need 2 functions: slack integration and report builder func (slack and report folders respectively). The API Gateway should attach to slack integration function. 
  
  The lambda function would need starting with PROFILE=<your-profile>. If you take a look config folder, your settings would be placed in config file under pattern: config-<your_profile>.json. For example, if you start your lambda function with PROFILE=test, the function would find settings from file config-test.json in config folder. 
  
  
  ### Slack configuration 

You need to create Slack Command in Slack App. 
  
  ![image](https://user-images.githubusercontent.com/10323118/124391524-32bc7880-dd1b-11eb-9f8b-3847ec5fb04f.png)

  ![image](https://user-images.githubusercontent.com/10323118/124391540-449e1b80-dd1b-11eb-98cd-647429becc17.png)

** Note: You need to enable "Escape channels, users, and links sent to your app" checkbox **   
  
  ![image](https://user-images.githubusercontent.com/10323118/124391574-71eac980-dd1b-11eb-9993-c840f1d7a00d.png)

## Usage: 
  
  ![image](https://user-images.githubusercontent.com/10323118/124391623-b4140b00-dd1b-11eb-80d7-78e1af8e711e.png)

  
  


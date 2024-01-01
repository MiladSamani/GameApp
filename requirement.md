# Game Application

## use case
### user use-cases
#### register
user can register to application by phone number
#### login
user can log in to application by phone number and password

## game uses-cases
#### Each game have a given number of questions
#### The difficulty level of questions are "easy, medium, hard"
#### Game winner is determined by number of correct answers that each user answered
#### Each game is belonged to specific category: sport, history, etc

## entity
### user
- ID 
- Phone number
- Avatar
- Name

### game
- ID
- Category
- Question List
- Players
- Winner

### question
- ID
- Question
- Answer List
- Correct Answer
- Difficulty
- Category
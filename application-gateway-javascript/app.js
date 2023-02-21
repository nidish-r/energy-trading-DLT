const path = require('path');

const express = require('express');
const bodyParser = require('body-parser');

const adminData = require('./routes/admin');
const overviewRoutes = require('./routes/overview');

const app = express();

app.set('view engine', 'ejs');

//  This helps in parsing text kind of data in json format
app.use(bodyParser.urlencoded({ extended: false }));

//  This helps to make public folder as accesable directory
//  Now we can use /css/man.css directl as it will be serched 
//  in public folder only
app.use(express.static(path.join(__dirname,'public')));

app.use(adminData.rouutes);

app.use((req, res, next) => {
    res.status(404).render('404');
});

app.listen(3000);

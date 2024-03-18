<?php

use Phalcon\Autoload\Loader;
//use Phalcon\Db\Adapter\Pdo\Mysql as DbAdapter;
use Phalcon\Di\FactoryDefault;
use Phalcon\Mvc\Application;
use Phalcon\Mvc\Url as UrlProvider;
use Phalcon\Mvc\View;
use Phalcon\Mvc\Router;

define('BASE_PATH', dirname(__DIR__));
const APP_PATH = BASE_PATH . '/src';

// Register an autoloader
$loader = new Loader();
$loader
    ->setDirectories(
        [
            APP_PATH . '/Controllers/',
        ]
    )
    ->setNamespaces(
        [
            'Comicshelf\\Controllers' => APP_PATH . '/Controllers/'
        ]
    )
    ->register()
;

// Create a DI
$container = new FactoryDefault();

/* @var Router $router */
$router = $container['router'];
$router->setDefaultNamespace('Comicshelf\\Controllers');

// Setting up the view component
$container['view'] = function () {
    $view = new View();
    $view->setViewsDir(APP_PATH . '/Views/');
    return $view;
};

// Setup a base URI so that all generated URIs include the "tutorial" folder
$container['url'] = function () {
    $url = new UrlProvider();
    $url->setBaseUri('/');
    return $url;
};

// Set the database service
//$container['db'] = function () {
//    return new DbAdapter(
//        [
//            "host"     => 'tutorial-mysql',
//            "username" => 'phalcon',
//            "password" => 'secret',
//            "dbname"   => 'phalcon_tutorial',
//        ]
//    );
//};

// Handle the request
try {
    $application = new Application($container);
    $response    = $application->handle($_SERVER["REQUEST_URI"]);
    $response->send();
} catch (Exception $e) {
    echo "Exception: ", $e->getMessage();
}
<?php

use Comicshelf\Services\MarvelClient;
use Phalcon\Assets\Manager;
use Phalcon\Autoload\Loader;
use Phalcon\Di\FactoryDefault;
use Phalcon\Logger\Adapter\Stream;
use Phalcon\Logger\Logger;
use Phalcon\Mvc\Application;
use Phalcon\Mvc\Router;
use Phalcon\Mvc\Url as UrlProvider;
use Phalcon\Mvc\View;
use Phalcon\Mvc\View\Engine\Volt;
use Phalcon\Mvc\ViewBaseInterface;

//use Phalcon\Db\Adapter\Pdo\Mysql as DbAdapter;

define('BASE_PATH', dirname(__DIR__));
const APP_PATH = BASE_PATH . '/src';

// Register an autoloader
$loader = new Loader();
$loader
    ->setDirectories(
        [
            APP_PATH . '/Controllers/',
            APP_PATH . '/Services/',
            APP_PATH . '/Util/',
        ]
    )
    ->setNamespaces(
        [
            'Comicshelf\\Controllers' => APP_PATH . '/Controllers/',
            'Comicshelf\\Services' => APP_PATH . '/Services/',
            'Comicshelf\\Util' => APP_PATH . '/Util/',
        ]
    )
    ->register()
;

// Create a DI
$container = new FactoryDefault();

$container['logger'] = function () {
    $stream = new Stream('php://stdout');
    return new Logger('messages', [
        'main' => $stream,
    ]);
};

$container["marvelClient"] = function () use ($container) {
    return new MarvelClient($container['logger']);
};

/* @var Router $router */
$router = $container['router'];
$router->setDefaultNamespace('Comicshelf\\Controllers');

// Setting up the view component

/** @var Manager $assets */
$assets = $container['assets'];
$assets->addCss('static/index.css');
$assets->addJs('static/index.js');

$container['voltService'] = function (ViewBaseInterface $view) use ($container) {
    $volt = new Volt($view, $container);
    $volt->setOptions(
        [
            'always' => true,
            'separator' => '_',
            'stat' => true,
            'path' => BASE_PATH . '/volt/',
            'stat' => true,
        ]
    );
    return $volt;
};

$container['view'] = function () {
    $view = new View();
    $view->setViewsDir(APP_PATH . '/Views/');
    $view->registerEngines(
        [
            '.phtml' => 'voltService',
        ]
    );
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
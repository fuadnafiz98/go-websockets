import { h } from 'preact';
import { Route, Router } from 'preact-router';
import Home from '../routes/home';

const App = () => (
	<div id="app">
		<main>
			<Router>
				<Route path="/" component={Home} />
			</Router>
		</main>
	</div>
);

export default App;

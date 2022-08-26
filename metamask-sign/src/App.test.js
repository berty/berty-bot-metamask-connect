import { render, screen } from '@testing-library/react';
import BaseApp from './pages/App';

test('renders learn react link', () => {
  render(<BaseApp />);
  const linkElement = screen.getByText(/learn react/i);
  expect(linkElement).toBeInTheDocument();
});

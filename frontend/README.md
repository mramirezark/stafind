# StaffFind Frontend - Next.js Application

This is the frontend application for StaffFind, built with Next.js 14, React 18, TypeScript, and Material-UI.

## 🚀 Features

- **Modern Stack**: Next.js 14 with App Router
- **TypeScript**: Full type safety
- **Material-UI**: Beautiful, responsive UI components
- **Real-time Updates**: Live reloading with Air for backend
- **Responsive Design**: Mobile-first approach
- **API Integration**: Seamless communication with Go backend

## 📁 Project Structure

```
frontend/
├── app/                          # Next.js App Router
│   ├── layout.tsx               # Root layout with theme
│   ├── page.tsx                 # Home page
│   ├── ThemeRegistry.tsx        # MUI theme provider
│   └── globals.css              # Global styles
├── components/                   # React components
│   ├── Dashboard.tsx            # Dashboard with stats
│   ├── EngineerManagement.tsx   # Engineer CRUD operations
│   ├── JobRequestForm.tsx       # Job request creation
│   ├── Navigation.tsx           # App navigation
│   └── SearchEngineers.tsx      # Engineer search
├── lib/                         # Utilities and configurations
│   └── api.ts                   # API client with Axios
├── next.config.js               # Next.js configuration
├── package.json                 # Dependencies and scripts
├── tsconfig.json                # TypeScript configuration
└── Dockerfile                   # Docker configuration
```

## 🛠️ Development Setup

### Prerequisites

- Node.js 18+
- npm or yarn
- Backend API running on port 8080

### Installation

```bash
# Install dependencies
npm install

# Start development server
npm run dev
```

The application will be available at `http://localhost:3000`.

### Available Scripts

```bash
npm run dev          # Start development server
npm run build        # Build for production
npm run start        # Start production server
npm run lint         # Run ESLint
npm run type-check   # Run TypeScript type checking
```

## 🎨 Components

### Dashboard
- Overview statistics (engineers, job requests, matches)
- Recent job requests table
- Real-time data from backend API

### Engineer Management
- Add new engineers with skills
- Edit existing engineer profiles
- View engineer list with skills
- Delete engineers

### Job Request Form
- Create new job requests
- Select required and preferred skills
- Set experience level and priority
- Department and location selection

### Search Engineers
- Search by skills (required/preferred)
- Filter by department and experience level
- Display engineer profiles with match scores
- Skill proficiency ratings

### Navigation
- Responsive navigation bar
- Mobile-friendly drawer
- Active state indicators
- Speed dial for quick actions

## 🔧 Configuration

### Environment Variables

Create a `.env.local` file:

```bash
NEXT_PUBLIC_API_URL=http://localhost:8080
NODE_ENV=development
```

### API Configuration

The API client is configured in `lib/api.ts`:
- Base URL from environment variables
- Request/response interceptors
- Error handling
- Authentication token support

### Theme Configuration

Material-UI theme is configured in `app/ThemeRegistry.tsx`:
- Primary/secondary colors
- Typography settings
- Component customizations
- Responsive breakpoints

## 🐳 Docker Support

### Development

```bash
# Using Makefile
make frontend-dev

# Or manually
cd frontend && npm run dev
```

### Production

```bash
# Build Docker image
make build-frontend

# Or manually
docker compose build frontend
```

## 📱 Responsive Design

The application is fully responsive with:
- Mobile-first design approach
- Breakpoint-specific layouts
- Touch-friendly interfaces
- Adaptive navigation (desktop menu vs mobile drawer)

## 🔗 API Integration

### Endpoints Used

- `GET /api/v1/engineers` - List engineers
- `POST /api/v1/engineers` - Create engineer
- `PUT /api/v1/engineers/:id` - Update engineer
- `DELETE /api/v1/engineers/:id` - Delete engineer
- `GET /api/v1/job-requests` - List job requests
- `POST /api/v1/job-requests` - Create job request
- `GET /api/v1/skills` - List skills
- `POST /api/v1/search` - Search engineers

### Data Flow

1. Components fetch data using SWR or direct API calls
2. Loading states are handled gracefully
3. Error states show user-friendly messages
4. Success states provide feedback

## 🎯 Key Features

### Real-time Updates
- Live reloading during development
- Automatic refresh on code changes
- Hot module replacement

### Type Safety
- Full TypeScript coverage
- Strict type checking
- Interface definitions for all data models

### Performance
- Next.js optimizations (code splitting, image optimization)
- Material-UI tree shaking
- Efficient re-rendering with React 18

### User Experience
- Smooth animations and transitions
- Consistent design language
- Intuitive navigation
- Accessible components

## 🚀 Deployment

### Production Build

```bash
npm run build
```

### Docker Production

```bash
docker compose up -d frontend
```

### Environment Configuration

For production deployment, ensure:
- `NEXT_PUBLIC_API_URL` points to production backend
- `NODE_ENV=production`
- Proper CORS configuration on backend

## 🔧 Troubleshooting

### Common Issues

1. **Port 3000 in use**
   ```bash
   lsof -ti:3000 | xargs kill -9
   ```

2. **API connection failed**
   - Check if backend is running on port 8080
   - Verify `NEXT_PUBLIC_API_URL` in environment

3. **Build errors**
   ```bash
   rm -rf .next node_modules
   npm install
   npm run build
   ```

4. **Type errors**
   ```bash
   npm run type-check
   ```

### Development Tips

- Use browser dev tools for debugging
- Check network tab for API calls
- Use React DevTools for component inspection
- Monitor console for errors and warnings

## 📚 Resources

- [Next.js Documentation](https://nextjs.org/docs)
- [Material-UI Documentation](https://mui.com/)
- [TypeScript Handbook](https://www.typescriptlang.org/docs/)
- [React 18 Features](https://react.dev/)

## 🤝 Contributing

1. Follow TypeScript best practices
2. Use Material-UI components consistently
3. Write responsive, accessible components
4. Test on multiple devices and browsers
5. Follow the established project structure
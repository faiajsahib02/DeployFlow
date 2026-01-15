import { useState } from 'react';
import { Dashboard } from '@/components/Dashboard';
import { DeploymentWorkspace } from '@/components/DeploymentWorkspace';
import { Project } from '@/services/api';

function App() {
  const [currentProject, setCurrentProject] = useState<Project | null>(null);

  if (currentProject) {
    return (
      <DeploymentWorkspace
        project={currentProject}
        onBack={() => setCurrentProject(null)}
      />
    );
  }

  return <Dashboard onProjectCreated={setCurrentProject} />;
}

export default App;

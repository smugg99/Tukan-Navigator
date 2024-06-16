<template>
  <v-container fluid class="fill-height">
    <v-row class="fill-height">
      <v-col class="fill-height">
        <!-- SVG container with relative positioning -->
        <div class="svg-container">
          <!-- SVG graph with absolute positioning -->
          <svg
            ref="svg"
            class="graph-svg"
            @mousedown="startInteraction"
            @mouseup="stopInteraction"
            @mouseleave="stopInteraction"
            @mousemove="handleMouseOver"
            @click="handleSvgClick"
          >
            <!-- Existing SVG content -->
            <g :transform="`translate(${panX}, ${panY})`">
              <!-- Render existing edges -->
              <Edge
                v-for="edge in validEdges"
                :key="'edge-' + edge.id"
                :from="findNode(edge.from)"
                :to="findNode(edge.to)"
                :edge="edge"
                :highlighted="(selectedEdgeId === edge.id || highlightedEdgeId === edge.id) && (mode === 'edit' || mode === 'remove')"
                @select="handleEdgeSelect(edge)"
              />
              <!-- Render existing nodes -->
              <Node
                v-for="node in nodes"
                :key="'node-' + node.id"
                :id="node.id"
                :x="node.x"
                :y="node.y"
                :selected="node.id === selectedNodeId || node.id === edgeStartNode"
                :highlighted="node.id === hoveredNodeId || node.id === selectedNodeIdInEditMode"
                :mode="mode"
                @select="selectNode"
                @mousedown.stop="startNodeDrag($event, node.id)"
              />
            </g>

            <!-- Ghost element for new edge -->
            <line
              v-if="mode === 'addEdge' && edgeStartNode && hoveredNodeId"
              :x1="findNode(edgeStartNode).x + panX"
              :y1="findNode(edgeStartNode).y + panY"
              :x2="findNode(hoveredNodeId).x + panX"
              :y2="findNode(hoveredNodeId).y + panY"
              stroke="black"
              opacity="0.5"
              stroke-dasharray="5,5"
              pointer-events="none"
            />

            <!-- Ghost element for new node -->
            <circle
              v-if="mode === 'add' && addNodeHovered"
              :cx="addNodeHovered.x + panX"
              :cy="addNodeHovered.y + panY"
              r="20"
              fill="transparent"
              stroke="black"
              opacity="0.5"
              stroke-dasharray="5,5"
              pointer-events="none"
            />

            <!-- Indicator for current position -->
            <circle
              v-if="isPathSet"
              :cx="currentPosition.x + panX"
              :cy="currentPosition.y + panY"
              r="10"
              fill="blue"
            />
          </svg>

          <!-- Button group outside SVG for z-index stacking -->
          <div class="button-group">
            <!-- Buttons here -->
            <v-btn
              :class="{ 'v-btn--active': mode === 'pan' }"
              @click.stop="setMode('pan')"
              dense
            >Move</v-btn>
            <v-btn
              :class="{ 'v-btn--active': mode === 'drag' }"
              @click.stop="setMode('drag')"
              dense
            >Drag</v-btn>
            <v-btn
              :class="{ 'v-btn--active': mode === 'add' }"
              @click.stop="setMode('add')"
              dense
            >Add Node</v-btn>
            <v-btn
              :class="{ 'v-btn--active': mode === 'remove' }"
              @click.stop="setMode('remove')"
              dense
            >Remove</v-btn>
            <v-btn
              :class="{ 'v-btn--active': mode === 'edit' }"
              @click.stop="setMode('edit')"
              dense
            >Edit</v-btn>
            <v-btn
              :class="{ 'v-btn--active': mode === 'addEdge' }"
              @click.stop="setMode('addEdge')"
              dense
            >Add Edge</v-btn>
            <v-btn @click.stop="panToNode('S')" dense>Go to start</v-btn>
            <v-btn
              :class="animationError ? 'error' : (animationRunning ? 'warning' : 'success')"
              @click.stop="toggleAnimation"
              dense
            >
              {{ animationRunning ? 'Stop' : (animationError ? 'Restart' : 'Start') }}
            </v-btn>
          </div>
        </div>
      </v-col>
    </v-row>
  </v-container>
</template>
<script>
import { throttle } from 'lodash';
import Node from './Node.vue';
import Edge from './Edge.vue';
import Toucan from './Toucan.vue';

export default {
  components: {
    Node,
    Edge,
    Toucan
  },
  data() {
    return {
      nodes: [
        { id: 'S', x: 50, y: 50 },
        { id: 'P', x: 250, y: 50 },
      ],
      edges: [],
      selectedNodeId: null,
      hoveredNodeId: null,
      selectedEdgeId: null,
      highlightedEdgeId: null,
      selectedNodeIdInEditMode: null,
      draggingNodeId: null,
      isPanning: false,
      offsetX: 0,
      offsetY: 0,
      panX: 0,
      panY: 0,
      panStartX: 0,
      panStartY: 0,
      width: 800,
      height: 600,
      mode: 'pan',
      creatingEdge: null,
      edgeStartNode: null,
      addNodeHovered: null,
      animationRunning: false,
      animationError: false,
      animatedPath: [],
      pathIndex: 0,
    };
  },
  computed: {
    validEdges() {
      return this.edges.filter(edge => this.findNode(edge.from) && this.findNode(edge.to));
    }
  },
  mounted() {
    this.throttledMouseMove = throttle(this.onMouseMove, 16);
    window.addEventListener('mousemove', this.throttledMouseMove);
    window.addEventListener('mouseup', this.stopInteraction);
  },
  beforeDestroy() {
    window.removeEventListener('mousemove', this.throttledMouseMove);
    window.removeEventListener('mouseup', this.stopInteraction);
  },
  methods: {
    findNode(id) {
      return this.nodes.find(node => node.id === id);
    },
    selectNode(id) {
      if (this.mode === 'remove') {
        this.removeNode(id);
      } else if (this.mode === 'edit') {
        if (id === 'S' || id === 'P') {
          alert('Node cannot be updated.');
          return;
        }

        this.selectedNodeIdInEditMode = id;
        const newId = prompt('Enter new ID', id);
        if (newId && newId !== id) {
          if (!this.findNode(newId)) {
            this.updateNodeId(id, newId);
          } else {
            alert('Duplicate ID detected. Please enter a unique ID.');
          }
        }
        this.selectedNodeIdInEditMode = null;
      } else if (this.mode === 'addEdge') {
        if (!this.edgeStartNode) {
          this.edgeStartNode = id;
          this.selectedNodeId = id;
        } else {
          if (this.edgeStartNode !== id) {
            const existingEdge = this.edges.find(edge => 
              (edge.from === this.edgeStartNode && edge.to === id) ||
              (edge.from === id && edge.to === this.edgeStartNode)
            );
            if (!existingEdge) {
              const weight = prompt('Enter weight for new edge');
              if (weight) {
                this.addEdge(this.edgeStartNode, id, parseInt(weight));
                this.edgeStartNode = null;
                this.selectedNodeId = null;
              }
            } else {
              alert('An edge already exists between these nodes.');
            }
          } else {
            alert('Cannot connect a node to itself. Select a different node.');
            this.edgeStartNode = null;
            this.selectedNodeId = null;
          }
        }
      } else if (this.mode === 'pan' || this.mode === 'drag') {
        this.selectedNodeId = id;
      }
    },
    handleEdgeSelect(edge) {
      if (this.mode === 'remove') {
        this.removeEdge(edge.id);
      } else if (this.mode === 'edit') {
        this.selectedEdgeId = edge.id;
        const newWeight = prompt('Enter new weight', edge.weight);
        if (newWeight) {
          this.updateEdgeWeight(edge.id, parseInt(newWeight));
        }
        this.selectedEdgeId = null;
      }
    },
    startNodeDrag(event, nodeId) {
      if (this.mode === 'drag') {
        event.preventDefault();
        this.draggingNodeId = nodeId;
        const node = this.findNode(nodeId);
        if (node) {
          this.offsetX = event.clientX - node.x - this.panX;
          this.offsetY = event.clientY - node.y - this.panY;
        }
      }
    },
    startInteraction(event) {
      if (this.mode === 'pan') {
        this.isPanning = true;
        this.panStartX = event.clientX;
        this.panStartY = event.clientY;
      }
      if (this.mode === 'drag' && this.draggingNodeId) {
        this.selectedNodeId = this.draggingNodeId;
      }
    },
    stopInteraction() {
      this.isPanning = false;
      this.draggingNodeId = null;
    },
    onMouseMove(event) {
      if (this.draggingNodeId !== null && this.mode === 'drag') {
        const node = this.findNode(this.draggingNodeId);
        if (node) {
          node.x = event.clientX - this.offsetX - this.panX;
          node.y = event.clientY - this.offsetY - this.panY;
        }
      }
      if (this.isPanning) {
        const dx = event.clientX - this.panStartX;
        const dy = event.clientY - this.panStartY;
        this.panX += dx;
        this.panY += dy;
        this.panStartX = event.clientX;
        this.panStartY = event.clientY;
      }
    },
    handleMouseOver(event) {
      const svgRect = this.$refs.svg.getBoundingClientRect();
      const x = event.clientX - svgRect.left - this.panX;
      const y = event.clientY - svgRect.top - this.panY;
      const hoveredNode = this.nodes.find(node => Math.abs(node.x - x) < 20 && Math.abs(node.y - y) < 20);

      if (this.mode === 'addEdge' || this.mode === 'edit' || this.mode === 'drag' || this.mode === 'remove') {
        this.hoveredNodeId = hoveredNode ? hoveredNode.id : null;
        if ((this.mode === 'edit' || this.mode === 'remove') && !hoveredNode) {
          const hoveredEdge = this.edges.find(edge => {
            const fromNode = this.findNode(edge.from);
            const toNode = this.findNode(edge.to);
            if (!fromNode || !toNode) return false;
            const dist = Math.abs(
              (toNode.y - fromNode.y) * x - (toNode.x - fromNode.x) * y + toNode.x * fromNode.y - toNode.y * fromNode.x
            ) /
              Math.sqrt(
                Math.pow(toNode.y - fromNode.y, 2) + Math.pow(toNode.x - fromNode.x, 2)
              );
            return dist < 5;
          });
          this.highlightedEdgeId = hoveredEdge ? hoveredEdge.id : null;
        } else {
          this.highlightedEdgeId = null;
        }
      } else if (this.mode === 'add') {
        this.addNodeHovered = { x, y };
      } else {
        this.hoveredNodeId = null;
        this.addNodeHovered = null;
        this.highlightedEdgeId = null;
      }
      if (this.draggingNodeId === null) {
        this.selectedNodeId = this.hoveredNodeId ? this.hoveredNodeId : null;
      }
    },
    handleSvgClick(event) {
      if (this.mode === 'add') {
        const svgRect = this.$refs.svg.getBoundingClientRect();
        const x = event.clientX - svgRect.left - this.panX;
        const y = event.clientY - svgRect.top - this.panY;
        const newId = prompt('Enter new ID');
        if (newId) {
          if (!this.findNode(newId)) {
            this.addNode(newId, x, y);
          } else {
            alert('Duplicate ID detected. Please enter a unique ID.');
          }
        }
      } else {
        if(this.mode !== 'addEdge') {
          this.selectedNodeId = null;
          this.edgeStartNode = null;
        }
      }
      this.addNodeHovered = null;
    },
    addNode(id, x, y) {
      if (!this.findNode(id)) {
        this.nodes.push({ id, x, y });
      } else {
        alert('Duplicate ID detected. Please enter a unique ID.');
      }
    },
    removeNode(id) {
      if (id === 'S' || id === 'P') {
        alert('Node cannot be deleted.');
        return;
      }

      this.nodes = this.nodes.filter(node => node.id !== id);
      this.edges = this.edges.filter(edge => edge.from !== id && edge.to !== id);
    },
    addEdge(from, to, weight) {
      if (!this.edges.find(edge => (edge.from === from && edge.to === to) || (edge.from === to && edge.to === from))) {
        this.edges.push({ id: `${from}-${to}`, from, to, weight });
      }
    },
    removeEdge(id) {
      this.edges = this.edges.filter(edge => edge.id !== id);
    },
    updateNodeId(oldId, newId) {
      const node = this.findNode(oldId);
      if (node) {
        node.id = newId;
        this.edges.forEach(edge => {
          if (edge.from === oldId) edge.from = newId;
          if (edge.to === oldId) edge.to = newId;
        });
      }
    },
    updateEdgeWeight(id, newWeight) {
      const edge = this.edges.find(edge => edge.id === id);
      if (edge) {
        edge.weight = newWeight;
      }
    },
    setMode(mode) {
      this.mode = mode;
      this.isPanning = false;
      this.draggingNodeId = null;
      this.edgeStartNode = null;
      this.selectedNodeId = null;
      this.hoveredNodeId = null;
      this.highlightedEdgeId = null;
      this.addNodeHovered = null;
    },
    getGraphRelations() {
      const nodes = this.nodes.map(node => node.id);
      const edges = this.edges.map(edge => ({
        from: edge.from,
        to: edge.to,
        weight: edge.weight
      }));

      const graphRelations = {
        nodes,
        edges
      };
      console.log(graphRelations);

      return graphRelations;
    },
    getGraphRelationsAndPositions() {
      const nodes = this.nodes.map(node => ({
        id: node.id,
        x: node.x,
        y: node.y
      }));

      const edges = this.edges.map(edge => ({
        from: edge.from,
        to: edge.to,
        weight: edge.weight
      }));

      const graphRelationsAndPositions = {
        nodes,
        edges
      };
      console.log(graphRelationsAndPositions);

      return graphRelationsAndPositions;
    },

    async toggleAnimation() {
      this.animationError = false;
      
      if (!this.animationRunning) {
        this.animationRunning = true;
        await this.startAnimation();
      } else {
        this.animationRunning = false;
      }
    },

    async startAnimation() {
      try {
        const response = await fetch('http://localhost:8000/api/v1/graphs/shortest-path', {
          method: 'POST',
          headers: {
            'Content-Type': 'application/json'
          },
          body: JSON.stringify({
            nodes: this.nodes.map(node => node.id),
            edges: this.edges.map(edge => ({
              from: edge.from,
              to: edge.to,
              weight: edge.weight
            }))
          })
        });
        
        const data = await response.json();
        if (data.path) {
          // Clear existing animation data
          this.animatedPath = [];

          // Start animation from the 'S' node
          this.animatedPath.push(this.findNode('S'));

          // Append the shortest path nodes to animatedPath
          this.animatedPath.push(...data.path.map(nodeId => this.findNode(nodeId)));
          this.pathIndex = 0;
          
          this.animationRunning = true;
          this.animationError = false;

          this.animateToucan();
        } else {
          console.error('Path data is not valid:', data);
          setTimeout(() => { this.animationRunning = false; this.animationError = true; }, 200);
        }
      } catch (error) {
        console.error('Error fetching path:', error);
        setTimeout(() => { this.animationRunning = false; this.animationError = true; }, 200);
      }
    },

    async animateToucan() {
      if (this.animationRunning && this.pathIndex < this.animatedPath.length) {
        const node = this.animatedPath[this.pathIndex];
        await this.panToucanToNode(node.id);
        this.pathIndex++;
        setTimeout(this.animateToucan, 1500);
      } else {
        this.animationRunning = false;
      }
    },

    async panToucanToNode(nodeId) {
      const node = this.findNode(nodeId);
      if (node) {
        const svgRect = this.$refs.svg.getBoundingClientRect();
        const centerX = svgRect.width / 2;
        const centerY = svgRect.height / 2;

        const currentPanX = this.panX;
        const currentPanY = this.panY;

        const targetX = -node.x + centerX - currentPanX;
        const targetY = -node.y + centerY - currentPanY;

        const duration = 500; // Duration in milliseconds
        let startTime = null;

        const animate = (currentTime) => {
          if (!startTime) startTime = currentTime;
          const elapsedTime = currentTime - startTime;
          const progress = Math.min(elapsedTime / duration, 1);

          const easingProgress = this.easeInOutQuad(progress);

          this.panX = currentPanX + targetX * easingProgress;
          this.panY = currentPanY + targetY * easingProgress;

          if (progress < 1) {
            requestAnimationFrame(animate);
          }
        };

        requestAnimationFrame(animate);
      }
    },

    async panToNode(nodeId) {
      const node = this.findNode(nodeId);
      if (node) {
        const svgRect = this.$refs.svg.getBoundingClientRect();
        const centerX = svgRect.width / 2;
        const centerY = svgRect.height / 2;

        const currentPanX = this.panX;
        const currentPanY = this.panY;

        const targetX = -node.x + centerX - currentPanX;
        const targetY = -node.y + centerY - currentPanY;

        const duration = 500;
        let startTime = null;

        const animate = (currentTime) => {
          if (!startTime) startTime = currentTime;
          const elapsedTime = currentTime - startTime;
          const progress = Math.min(elapsedTime / duration, 1);

          const easingProgress = this.easeInOutQuad(progress);

          this.panX = currentPanX + targetX * easingProgress;
          this.panY = currentPanY + targetY * easingProgress;

          if (progress < 1) {
            requestAnimationFrame(animate);
          }
        };

        requestAnimationFrame(animate);
      }
    },

    easeInOutQuad(t) {
      return t < 0.5 ? 2 * t * t : -1 + (4 - 2 * t) * t;
    },
  }
};
</script>

<style scoped>
.fill-height {
  height: 100%;
  margin: 0;
  padding: 0;
  overflow: hidden;
}
.svg-container {
  position: relative;
  width: 100%;
  height: 100%;
  overflow: hidden;
}
.graph-svg {
  position: absolute;
  top: 0;
  left: 0;
  width: 100%;
  height: 100%;
  z-index: 1; /* Ensure SVG content is below buttons */
}
.button-group {
  user-select: none;
  position: absolute;
  top: 10px; /* Adjust as needed */
  left: 10px; /* Adjust as needed */
  z-index: 2; /* Ensure buttons are above SVG content */
}
</style>
